package service_internal

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/shared/models/pb/access"
)

func (s *ServiceStruct) CreateJwt(ctx context.Context, exp time.Duration, v *models.ValidatedCodeCheckReq) (string, error) {
	jti := uuid.New()
	s.key.Jti = jti // до конвертации в строку, чтоб потом записать в базу с типом uuid
	s.key.Alg = "RS256"

	jtiStr := jti.String()

	payload := jwt.MapClaims{
		"jti":      jtiStr,
		"version":  1,
		"shard_id": rand.Intn(9) + 1, // на всякий случай указываем на будущее чтоб иметь вечные токены разделенные на 9 частей
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token.Header["kid"] = s.key.PubkeyId

	tokenStr, err := token.SignedString(s.key.PrivateKey)
	if err != nil {
		return "", err
	}

	// сначала добавляем запись в whitelist, потому пишем в базу, так как вызов внешний и рисков больше, чтоб точно быть уверенными что если в postgres есть запись, то в вайтлист jti добавлен
	// TODO реализовать retry для временных ошибок
	resp, err := s.accessClient.AddWhitelist(ctx, &access.SendRequest{Jti: jtiStr})
	if err != nil {
		return "", fmt.Errorf("addWhitelist rpc failed: %w", err)
	}
	if resp == nil {
		return "", errors.New("empty server response")
	}

	// GetSuccess это сгенерированный protoc'ом геттер, обязательно используем геттеры при работы с grpc!
	if !resp.GetSuccess() {
		return "", errors.New("jti not whitelisted")
	}

	if err = s.repo.AddSignedToken(ctx, v, s.key); err != nil {
		return "", err
	}

	return tokenStr, nil

}
