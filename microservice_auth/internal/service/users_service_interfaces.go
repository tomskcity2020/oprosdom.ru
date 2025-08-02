package service

import (
	"context"

	"oprosdom.ru/microservice_auth/internal/biz"
	"oprosdom.ru/microservice_auth/internal/models"
	"oprosdom.ru/microservice_auth/internal/repo"
	service_internal "oprosdom.ru/microservice_auth/internal/service/internal"
	"oprosdom.ru/microservice_auth/internal/transport"
)

type UsersService interface {
	PhoneSend(ctx context.Context, p *models.ValidatedPhoneSendReq) error
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(repo repo.RepositoryInterface, codeTransport transport.TransportInterface) UsersService {
	biz := biz.NewBizFactory()
	return service_internal.NewCallInternalService(repo, biz, codeTransport)
}
