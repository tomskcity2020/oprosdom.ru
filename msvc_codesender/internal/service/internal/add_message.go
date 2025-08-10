package service_internal

import (
	"context"

	"oprosdom.ru/msvc_codesender/internal/models"
)

func (s *ServiceStruct) AddMessage(ctx context.Context, validMsg *models.ValidatedMsg) error {

	// проверяем  type  номера: mobile в sms_, landline и unknown в calls_
	switch validMsg.Type {
	case "mobile":
		s.repo.InsertSms(ctx, validMsg)
	case "landline", "unknown":
		s.repo.InsertCall(ctx, validMsg)
	}
	return nil

}
