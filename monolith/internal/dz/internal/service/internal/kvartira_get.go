package service_internal

import (
	"errors"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartiraGet(id string) (*models.Kvartira, error) {

	// чекаем что id является корректным uuid
	if err := s.biz.UuidCheck(id); err != nil {
		return nil, errors.New("неправильный id квартиры")
	}

	data, err := s.repo.GetKvartiraById(id)
	if err != nil {
		if err.Error() == "not_found" {
			return nil, errors.New("квартира не найдена")
		} else {
			return nil, errors.New("сбой данные по квартире не получены")
		}
	}

	return data, nil

}
