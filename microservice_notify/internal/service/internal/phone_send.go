package service_internal

import (
	"context"

	"oprosdom.ru/microservice_notify/internal/models"
)

func (s *ServiceStruct) ProcessMessage(ctx context.Context, validMsg *models.ValidatedMsg) error {

	// записываем в репо то что получили из кафки только добавляем доп столбцы

	return nil

}
