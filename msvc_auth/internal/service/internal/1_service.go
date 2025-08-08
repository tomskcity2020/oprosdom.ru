package service_internal

import (
	"oprosdom.ru/msvc_auth/internal/biz"
	"oprosdom.ru/msvc_auth/internal/repo"
	"oprosdom.ru/msvc_auth/internal/transport"
)

type ServiceStruct struct {
	ramRepo       repo.RamRepoInterface
	repo          repo.RepositoryInterface
	biz           biz.BizInterface
	codeTransport transport.TransportInterface
}

func NewCallInternalService(ramRepo repo.RamRepoInterface, repo repo.RepositoryInterface, biz biz.BizInterface, codeTransport transport.TransportInterface) *ServiceStruct {
	return &ServiceStruct{
		ramRepo:       ramRepo,
		repo:          repo,
		biz:           biz,
		codeTransport: codeTransport,
	}
}
