package service_internal

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"oprosdom.ru/msvc_auth/internal/models"
)

func (s *ServiceStruct) CreateJwt(ctx context.Context, exp time.Duration, v *models.ValidatedCodeCheckReq) (string, error) {
	jti := uuid.New()
	s.key.Jti = jti // до конвертации в строку, чтоб потом записать в базу с типом uuid
	s.key.Alg = "RS256"

	payload := jwt.MapClaims{
		"jti": jti.String(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token.Header["kid"] = s.key.PubkeyId

	tokenStr, err := token.SignedString(s.key.PrivateKey)
	if err != nil {
		return "", err
	}

	s.repo.AddSignedToken(ctx, v, s.key)

	return tokenStr, nil

}
