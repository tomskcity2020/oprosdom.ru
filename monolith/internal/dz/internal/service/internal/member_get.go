package service_internal

import (
	"context"
	"errors"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) MemberGet(ctx context.Context, id string) (*models.Member, error) {

	// чекаем что id является корректным uuid
	if err := s.biz.UuidCheck(id); err != nil {
		return nil, errors.New("неправильный id жителя")
	}

	m, err := s.repo.MemberGetById(ctx, id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return m, nil

}
