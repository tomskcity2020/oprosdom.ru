package users_service_internal

import (
	users_biz "oprosdom.ru/core/internal/users/biz"
	users_repo "oprosdom.ru/core/internal/users/repo"
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
