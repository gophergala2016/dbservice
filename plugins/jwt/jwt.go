package jwt

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/gophergala2016/dbserver/plugins"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type JWT struct {
	Secret           string
	Issuer           string
	ExpirationTime   duration `toml:"expiration"`
	RotationDeadline duration `toml:"rotation_deadline"`
}

func (self *JWT) ParseConfig(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading plugin config: %v", err))
	}
	_, err = toml.Decode(string(content), self)
	return err
}

func (self *JWT) Process(data map[string]interface{}, arg map[string]interface{}) *plugins.Response {
	response := &plugins.Response{}
	if data["__jwt"] == nil {
		response.Data = data
		return response
	}
	payload, ok := data["__jwt"].(map[string]interface{})
	if !ok {
		response.ResponseCode = 500
		response.Error = fmt.Sprintf("__jwt parameter doesn't contain hash, but %v", data["__jwt"])
		return response
	}
	delete(data, "__jwt")
	response.Data = data
	token, err := self.GenerateToken(payload)
	if err != nil {
		response.ResponseCode = 500
		response.Error = err.Error()
		return response
	}
	if len(token) > 0 {
		response.Headers = make(map[string][]string)
		response.Headers["Authorization"] = []string{"Bearer " + string(token)}
	}
	return response
}

func (self *JWT) GenerateToken(payload map[string]interface{}) ([]byte, error) {
	claims := jws.Claims{}
	for key, value := range payload {
		claims.Set(key, value)
	}
	if self.Issuer != "" {
		claims.SetIssuer(self.Issuer)
	}
	if self.ExpirationTime.Duration > 0 {
		claims.SetExpiration(time.Now().Add(self.ExpirationTime.Duration))
	}
	token := jws.NewJWT(claims, crypto.SigningMethodHS256)
	serializedToken, err := token.Serialize([]byte(self.Secret))
	if err != nil {
		return nil, err
	}
	return serializedToken, nil
}

func (self *JWT) ProcessBeforeHook(data map[string]interface{}, r *http.Request) *plugins.Response {
	headerValue := r.Header.Get("Authorization")
	if headerValue == "" {
		return nil
	}
	if !strings.HasPrefix(headerValue, "Bearer ") {
		return nil
	}
	headerValue = strings.Replace(headerValue, "Bearer ", "", 1)
	token, err := jws.ParseJWT([]byte(headerValue))
	if err != nil {
		return nil
	}
	err = token.Validate([]byte(self.Secret), crypto.SigningMethodHS256)
	if err != nil {
		return nil
	}
	expiration, ok := token.Claims().Expiration()
	if !ok {
		return nil
	}
	if expiration.Unix() < time.Now().Unix() {
		return nil
	}
	if time.Now().Add(self.RotationDeadline.Duration).Unix() > expiration.Unix() {
		token, err := self.GenerateToken(token.Claims())
		response := &plugins.Response{}
		if err != nil {
			response.ResponseCode = 500
			response.Error = err.Error()
			return response
		}
		if len(token) > 0 {
			response.Headers = make(map[string][]string)
			response.Headers["Authorization"] = []string{"Bearer " + string(token)}
		}
	}
	data["jwt"] = token.Claims()
	return nil
}
