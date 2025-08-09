package service_internal

import (
	"context"
	"errors"

	"oprosdom.ru/core/internal/dz/internal/models"
)

func (s *ServiceStruct) MembersGet(ctx context.Context) ([]*models.Member, error) {
	m, err := s.repo.MembersGet(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return m, nil
}
