package service

import (
	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	service_internal "oprosdom.ru/monolith/internal/dz/internal/service/internal"
)

type ServiceInterface interface {
	CountData()
	RunParallel(modelsData []models.ModelInterface)
	RunSeq(modelsData []models.ModelInterface)
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory() ServiceInterface {
	repository, err := repo.GetRepoSingleton()
	if err != nil {
		panic(err)
	}
	return service_internal.NewCallInternalService(repository)
}
