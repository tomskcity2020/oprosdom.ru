package service_internal

import (
	"errors"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) MemberGet(id string) (*models.Member, error) {

	// чекаем что id является корректным uuid
	if err := s.biz.UuidCheck(id); err != nil {
		return nil, errors.New("неправильный id жителя")
	}

	data, err := s.repo.GetMemberById(id)
	if err != nil {
		if err.Error() == "not_found" {
			return nil, errors.New("житель не найден")
		} else {
			return nil, errors.New("данные по жителю не получены")
		}
	}

	return data, nil

}
