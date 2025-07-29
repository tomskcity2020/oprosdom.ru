package users_service

import (
	"context"

	users_biz "oprosdom.ru/monolith/internal/users/biz"
	users_models "oprosdom.ru/monolith/internal/users/models"
	users_repo "oprosdom.ru/monolith/internal/users/repo"
	users_service_internal "oprosdom.ru/monolith/internal/users/service/internal"
)

type UsersService interface {
	PhoneSend(ctx context.Context, p *users_models.ValidatedPhoneSendReq) error
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(repo users_repo.RepositoryInterface) UsersService {
	biz := users_biz.NewBizFactory()
	return users_service_internal.NewCallInternalService(repo, biz)
}
