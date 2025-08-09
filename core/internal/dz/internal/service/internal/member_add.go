package service_internal

import (
	"context"
	"errors"
	"fmt"

	"oprosdom.ru/core/internal/dz/internal/models"
)

func (s *ServiceStruct) MemberAdd(ctx context.Context, m *models.Member) error {

	// проводим первичную проверку (структура заполнена полностью сейчас)
	if err := s.biz.BasicMemberValidation(m); err != nil {
		return fmt.Errorf("basic member validation failed: %v", err.Error())
	}

	if err := s.repo.MemberAdd(ctx, m); err != nil {
		return errors.New(err.Error())
	}

	return nil

}
