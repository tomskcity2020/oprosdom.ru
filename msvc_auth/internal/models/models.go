package models

import "crypto/rsa"

type PhoneCode struct {
	Phone string `json:"phone"`
	Code  uint32 `json:"code"`
}

type KeyData struct {
	PrivateKey *rsa.PrivateKey
	PubkeyId   string
}
