package service_internal

import (
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartiraAdd(kvartira *models.Kvartira) error {

	// создаем uuid (отдельно делаем чтоб BasicValidate использовать и для POST и для PUT)
	id, err := s.biz.UuidCreate()
	if err != nil {
		return errors.New("create uuid failed")
	}

	kvartira.Id = id

	// проводим первичную проверку (структура заполнена полностью сейчас)
	if err := s.biz.BasicKvartiraValidation(kvartira); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	// сохраняем в файл
	if err := s.repo.SaveToFile(kvartira); err != nil {
		return errors.New("save to file failed")
	}

	// сохраняем в слайс
	if err := s.repo.Save(kvartira); err != nil {
		return errors.New("save to slice failed")
	}

	return nil

}
