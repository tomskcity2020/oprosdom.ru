package repo

import (
	"context"

	"oprosdom.ru/microservice_auth/internal/models"
	repo_internal "oprosdom.ru/microservice_auth/internal/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	PhoneSend(ctx context.Context, v *models.ValidatedPhoneSendReq) error
}
