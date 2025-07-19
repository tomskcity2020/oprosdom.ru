package service_internal

import (
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) MemberUpdate(member *models.Member) error {

	// проводим первичную валидацию (когда структура заполнена включая id)
	if err := s.biz.BasicMemberValidation(member); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	if err := s.repo.UpdateFile(member); err != nil {
		return errors.New("update file failed")
	}

	if err := s.repo.UpdateSlice(member); err != nil {
		return errors.New("update slice failed")
	}

	return nil

}
