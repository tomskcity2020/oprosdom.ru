package service_internal

import (
	"oprosdom.ru/core/internal/polls/biz"
	"oprosdom.ru/core/internal/polls/repo"
)

type ServiceStruct struct {
	repo polls_repo.RepositoryInterface
	biz  polls_biz.BizInterface
}

func NewCallInternalService(repo polls_repo.RepositoryInterface, biz polls_biz.BizInterface) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
		biz:  biz,
	}
}
