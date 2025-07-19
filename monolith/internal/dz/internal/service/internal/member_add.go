package service_internal

import (
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) MemberAdd(member *models.Member) error {

	id, err := s.biz.UuidCreate()
	if err != nil {
		return errors.New("create uuid failed")
	}

	member.Id = id

	// проводим первичную проверку (структура заполнена полностью сейчас)
	// передаем member, а не &member потому что member уже является указателем
	if err := s.biz.BasicMemberValidation(member); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	// сохраняем в файл
	if err := s.repo.SaveToFile(member); err != nil {
		return errors.New("save to file failed")
	}

	// сохраняем в слайс
	if err := s.repo.Save(member); err != nil {
		return errors.New("save to slice failed")
	}

	return nil
}
