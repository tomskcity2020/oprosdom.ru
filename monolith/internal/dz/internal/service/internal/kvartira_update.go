package service_internal

import (
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartiraUpdate(kvartira *models.Kvartira) error {

	// проводим первичную валидацию (когда структура заполнена включая id)
	if err := s.biz.BasicKvartiraValidation(kvartira); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	if err := s.repo.UpdateFile(kvartira); err != nil {
		return errors.New("update file failed")
	}

	if err := s.repo.UpdateSlice(kvartira); err != nil {
		return errors.New("update slice failed")
	}

	return nil

}
