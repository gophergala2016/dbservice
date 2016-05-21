package jwt

import (
	"time"
)

type JWT struct {
	Secret           []byte
	Issuer           string
	ExpirationTime   time.Time
	RotationDeadline time.Time
}

func (self *JWT) ParseConfig(path string) error {
	return nil
}

// app.Register(&JWT{}, "jwt") || app.Register("jwt", JWT)

// Hooks: 1. Before request
//        2. Process - when called in pipeline
