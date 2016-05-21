package jwt

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	jwtPlugin := &JWT{}
	err := jwtPlugin.ParseConfig("test_config/jwt.toml")
	if err != nil {
		t.Error(err)
	}
	if string(jwtPlugin.Secret) != "secret123" {
		t.Errorf("Secret key was parsed incorrectly. Expected: '%v', but got: '%v'", "secret123", string(jwtPlugin.Secret))
	}
}
