package service_internal

import (
	"context"
	"errors"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartiraGet(ctx context.Context, id string) (*models.Kvartira, error) {

	// чекаем что id является корректным uuid
	if err := s.biz.UuidCheck(id); err != nil {
		return nil, errors.New("неправильный id квартиры")
	}

	k, err := s.repo.KvartiraGetById(ctx, id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return k, nil

}
