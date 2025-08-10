package users_repo

import (
	"context"

	//"oprosdom.ru/core/internal/users/models"
	users_models "oprosdom.ru/core/internal/users/models"
	users_repo_internal "oprosdom.ru/core/internal/users/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return users_repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	PhoneSend(ctx context.Context, v *users_models.ValidatedPhoneSendReq) error
}
