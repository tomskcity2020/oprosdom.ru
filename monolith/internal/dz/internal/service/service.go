package service

import (
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	service_internal "oprosdom.ru/monolith/internal/dz/internal/service/internal"
)

type ServiceInterface interface {
	Run()
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory() ServiceInterface {
	return service_internal.NewCallInternalService(repo.NewRepoFactory())
}
