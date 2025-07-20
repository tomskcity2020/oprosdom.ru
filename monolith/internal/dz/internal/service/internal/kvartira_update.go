package service_internal

import (
	"context"
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartiraUpdate(ctx context.Context, k *models.Kvartira) error {

	// проводим первичную валидацию (когда структура заполнена включая id)
	// сначала отдельно проверяем id, так как из BasicValidation эта проверка исключена (ввиду генерации id на стороне базы данных)
	if err := s.biz.UuidCheck(k.Id); err != nil {
		return fmt.Errorf("id validation failed: %v", err.Error())
	}

	if err := s.biz.BasicKvartiraValidation(k); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	if err := s.repo.KvartiraUpdate(ctx, k); err != nil {
		return errors.New(err.Error())
	}

	return nil

}
