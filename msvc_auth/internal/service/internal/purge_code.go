package service_internal

import (
	"context"

	"oprosdom.ru/msvc_auth/internal/models"
)

func (s *ServiceStruct) PurgeCode(ctx context.Context, p *models.ValidatedCodeCheckReq) error {

	if _, err := s.ramRepo.Del(ctx, "phone:"+p.Phone); err != nil {
		return err
	}
	return nil
}
