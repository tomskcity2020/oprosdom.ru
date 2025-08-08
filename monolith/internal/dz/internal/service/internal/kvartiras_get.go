package service_internal

import (
	"context"
	"errors"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartirasGet(ctx context.Context) ([]*models.Kvartira, error) {
	k, err := s.repo.KvartirasGet(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return k, nil
}
