package service_internal

import (
	"context"
	"errors"
)

func (s *ServiceStruct) RemoveById(ctx context.Context, id string, mk string) error {

	// чекаем что id является корректным uuid
	if err := s.biz.UuidCheck(id); err != nil {
		return errors.New(err.Error())
	}

	if err := s.repo.DeleteById(ctx, id, mk); err != nil {
		return errors.New(err.Error())
	}

	return nil

}
