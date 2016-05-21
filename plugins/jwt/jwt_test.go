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
