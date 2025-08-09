package service

import (
	"context"
	"time"

	"oprosdom.ru/msvc_auth/internal/biz"
	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/msvc_auth/internal/repo"
	service_internal "oprosdom.ru/msvc_auth/internal/service/internal"
	"oprosdom.ru/msvc_auth/internal/transport"
)

type UsersService interface {
	PhoneSend(ctx context.Context, p *models.ValidatedPhoneSendReq) error
	CodeCheck(ctx context.Context, p *models.ValidatedCodeCheckReq) error
	CreateJwt(ctx context.Context, exp time.Duration, v *models.ValidatedCodeCheckReq) (string, error)
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(key *models.KeyData, ramRepo repo.RamRepoInterface, repo repo.RepositoryInterface, codeTransport transport.TransportInterface) UsersService {
	biz := biz.NewBizFactory()
	return service_internal.NewCallInternalService(key, ramRepo, repo, biz, codeTransport)
}
