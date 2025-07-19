package service_internal

import (
	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) MembersGet() ([]*models.Member, error) {
	// TODO добавить ошибку в репо и тут обрабатывать
	data := s.repo.GetSliceMembers()
	return data, nil
}
