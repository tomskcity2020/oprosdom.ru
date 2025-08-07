package service

import (
	"context"

	"oprosdom.ru/msvc_auth/internal/biz"
	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/msvc_auth/internal/repo"
	service_internal "oprosdom.ru/msvc_auth/internal/service/internal"
	"oprosdom.ru/msvc_auth/internal/transport"
)

type UsersService interface {
	PhoneSend(ctx context.Context, p *models.ValidatedPhoneSendReq) error
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(ramRepo repo.RamRepoInterface, repo repo.RepositoryInterface, codeTransport transport.TransportInterface) UsersService {
	biz := biz.NewBizFactory()
	return service_internal.NewCallInternalService(ramRepo, repo, biz, codeTransport)
}
