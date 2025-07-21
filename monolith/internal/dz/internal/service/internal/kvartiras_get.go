package service_internal

import (
	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) KvartirasGet() ([]*models.Kvartira, error) {
	// TODO добавить ошибку в репо и тут обрабатывать
	data := s.repo.GetSliceKvartiras()
	return data, nil
}
