package service_internal

import (
	"context"
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) MemberUpdate(ctx context.Context, m *models.Member) error {

	// проводим первичную валидацию (когда структура заполнена включая id)
	// сначала отдельно проверяем id, так как из BasicValidation эта проверка исключена (ввиду генерации id на стороне базы данных)
	if err := s.biz.UuidCheck(m.Id); err != nil {
		return fmt.Errorf("id validation failed: %v", err.Error())
	}

	if err := s.biz.BasicMemberValidation(m); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	if err := s.repo.MemberUpdate(ctx, m); err != nil {
		return errors.New(err.Error())
	}

	return nil

}
