package users_repo

import (
	"context"

	//"oprosdom.ru/monolith/internal/users/models"
	users_models "oprosdom.ru/monolith/internal/users/models"
	users_repo_internal "oprosdom.ru/monolith/internal/users/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return users_repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	PhoneSend(ctx context.Context, v *users_models.ValidatedPhoneSendReq) error
}
