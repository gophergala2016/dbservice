package jwt

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gophergala2016/dbserver/plugins"
	"io/ioutil"
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

func (self *JWT) Process(data map[string]interface{}) *plugins.Response {
	return nil
}

// app.Register(&JWT{}, "jwt") || app.Register("jwt", JWT)

// Hooks: 1. Before request
//        2. Process - when called in pipeline
