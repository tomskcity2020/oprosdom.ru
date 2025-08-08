package users_service_internal

import (
	users_biz "oprosdom.ru/monolith/internal/users/biz"
	users_repo "oprosdom.ru/monolith/internal/users/repo"
)

type ServiceStruct struct {
	repo users_repo.RepositoryInterface
	biz  users_biz.BizInterface
}

func NewCallInternalService(repo users_repo.RepositoryInterface, biz users_biz.BizInterface) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
		biz:  biz,
	}
}
