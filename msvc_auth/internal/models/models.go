package models

import (
	"crypto/rsa"

	"github.com/google/uuid"
)

type PhoneCode struct {
	Phone string `json:"phone"`
	Code  uint32 `json:"code"`
}

type KeyData struct {
	PrivateKey *rsa.PrivateKey
	PubkeyId   string
	Jti        uuid.UUID
	Alg        string
}
