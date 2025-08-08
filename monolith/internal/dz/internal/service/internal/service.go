package service_internal

import (
	"oprosdom.ru/monolith/internal/dz/internal/biz"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
)

type ServiceStruct struct {
	repo repo.RepositoryInterface
	biz  biz.BizInterface
}

func NewCallInternalService(repo repo.RepositoryInterface, biz biz.BizInterface) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
		biz:  biz,
	}
}
