package service

import (
	"context"

	"oprosdom.ru/msvc_codesender/internal/gateway"
	"oprosdom.ru/msvc_codesender/internal/models"
	"oprosdom.ru/msvc_codesender/internal/repo"
	service_internal "oprosdom.ru/msvc_codesender/internal/service/internal"
)

type ServiceInterface interface {
	AddMessage(ctx context.Context, validMsg *models.ValidatedMsg) error
	ProcessMessage(ctx context.Context, gateway *gateway.Gateway) error
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(repo repo.RepositoryInterface) ServiceInterface {
	return service_internal.NewCallInternalService(repo)
}
