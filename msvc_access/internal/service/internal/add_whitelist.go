package service_internal

import (
	"context"

	"oprosdom.ru/shared/models/pb/access"
)

func (s *ServiceStruct) AddWhitelist(ctx context.Context, req *access.SendRequest) (*access.SendResponse, error) {

	// к proto полям структур обращаемся только через геттеры иначе в теории можно столкнуться с nil'ом
	jti := req.GetJti()
	if jti == "" {
		return &access.SendResponse{Success: false}, nil
	}

	// префикс whitelist не указываем, экономим место так как whitelist будет расти постоянно и нужно быстро делать снэпшоты и бэкапы + redis только для whitelist будет использоваться
	if err := s.ramRepo.Set(ctx, jti, "0", 0); err != nil {
		return &access.SendResponse{Success: false}, nil
	}

	return &access.SendResponse{Success: true}, nil
}
