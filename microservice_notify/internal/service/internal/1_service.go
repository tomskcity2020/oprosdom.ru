package service_internal

import (
	"oprosdom.ru/microservice_notify/internal/repo"
)

type ServiceStruct struct {
	repo repo.RepositoryInterface
}

func NewCallInternalService(repo repo.RepositoryInterface) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
	}
}
