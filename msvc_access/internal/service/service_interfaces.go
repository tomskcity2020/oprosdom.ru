package service

import (
	"context"

	"oprosdom.ru/msvc_access/internal/repo"
	service_internal "oprosdom.ru/msvc_access/internal/service/internal"
	"oprosdom.ru/shared/models/pb/access"
)

type Service interface {
	AddWhitelist(ctx context.Context, req *access.SendRequest) (*access.SendResponse, error)
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(ramRepo repo.RamRepoInterface) Service {
	return service_internal.NewCallInternalService(ramRepo)
}
