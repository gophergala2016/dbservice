package jwt

import (
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	jwtPlugin := &JWT{}
	err := jwtPlugin.ParseConfig("test_config/jwt.toml")
	if err != nil {
		t.Error(err)
	}
	if jwtPlugin.Secret != "secret123" {
		t.Errorf("Secret attribute was parsed incorrectly. Expected: '%v', but got: '%v'", "secret123", jwtPlugin.Secret)
	}
	if jwtPlugin.Issuer != "issuer" {
		t.Errorf("Issuer attribute was parsed incorrectly. Expected: '%v', but got: '%v'", "issuer", jwtPlugin.Issuer)
	}
	if jwtPlugin.ExpirationTime.Duration != time.Duration(time.Hour*4) {
		t.Errorf("Expected 4 hour expiratio time, but got: %v", jwtPlugin.ExpirationTime.Duration)
	}
	if jwtPlugin.RotationDeadline.Duration != time.Duration(time.Hour*2) {
		t.Errorf("Expected 2 hour rotation deadline, but got: %v", jwtPlugin.RotationDeadline.Duration)
	}
}

func TestParseMultilineSecret(t *testing.T) {
	jwtPlugin := &JWT{}
	err := jwtPlugin.ParseConfig("test_config/jwt_multiline.toml")
	if err != nil {
		t.Error(err)
	}
	if jwtPlugin.Secret == "" {
		t.Errorf("Secret attribute was parsed incorrectly. Expected multiline, but got: '%v'", jwtPlugin.Secret)
	}
}

func TestProcess(t *testing.T) {
	jwt := &JWT{Secret: "secret"}
	data := make(map[string]interface{})
	data["__jwt"] = map[string]interface{}{"success": true}
	data["response"] = "data"
	response := jwt.Process(data, nil)
	if response == nil {
		t.Error("Not expected to get nil response after jwt processing")
	}
	if len(response.Headers["Authorization"]) != 1 {
		t.Error("Expected authorization header to be added by jwt plugin, but found none")
	}
	if response.Data["response"] != "data" {
		t.Error("Not found expected data after jwt plugin processing")
	}
	if response.Data["__jwt"] != nil {
		t.Error("__jwt token was supposed to be removed from data after jwt processing")
	}
}
