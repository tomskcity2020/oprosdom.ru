package service

import (
	"context"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	service_internal "oprosdom.ru/monolith/internal/dz/internal/service/internal"
)

type ServiceInterface interface {
	RunParallel(ctx context.Context, modelsData []models.ModelInterface)
	RunSeq(ctx context.Context, modelsData []models.ModelInterface)
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory() ServiceInterface {
	return service_internal.NewCallInternalService(repo.NewRepoFactory())
}
