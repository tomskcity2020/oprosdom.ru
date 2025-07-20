package service_internal

import (
	"context"
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartiraAdd(ctx context.Context, kvartira *models.Kvartira) error {

	// проводим первичную проверку (структура заполнена полностью сейчас)
	if err := s.biz.BasicKvartiraValidation(kvartira); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	if err := s.repo.KvartiraAdd(ctx, kvartira); err != nil {
		return errors.New(err.Error())
	}

	return nil

}
