package service_internal

import (
	"oprosdom.ru/msvc_auth/internal/biz"
	"oprosdom.ru/msvc_auth/internal/repo"
	"oprosdom.ru/msvc_auth/internal/transport"
)

type ServiceStruct struct {
	repo          repo.RepositoryInterface
	biz           biz.BizInterface
	codeTransport transport.TransportInterface
}

func NewCallInternalService(repo repo.RepositoryInterface, biz biz.BizInterface, codeTransport transport.TransportInterface) *ServiceStruct {
	return &ServiceStruct{
		repo:          repo,
		biz:           biz,
		codeTransport: codeTransport,
	}
}
