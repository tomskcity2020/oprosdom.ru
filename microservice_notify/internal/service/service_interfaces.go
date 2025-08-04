package service

import (
	"context"

	"oprosdom.ru/microservice_notify/internal/models"
	"oprosdom.ru/microservice_notify/internal/repo"
	service_internal "oprosdom.ru/microservice_notify/internal/service/internal"
)

type ServiceInterface interface {
	ProcessMessage(ctx context.Context, validMsg *models.ValidatedMsg) error
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(repo repo.RepositoryInterface) ServiceInterface {
	return service_internal.NewCallInternalService(repo)
}
