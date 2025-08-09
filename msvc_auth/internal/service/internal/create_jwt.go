package service_internal

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *ServiceStruct) CreateJwt(exp time.Duration) (string, error) {
	jti := uuid.NewString()

	payload := jwt.MapClaims{
		"jti": jti,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token.Header["kid"] = s.key.PubkeyId

	tokenStr, err := token.SignedString(s.key.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil

}
