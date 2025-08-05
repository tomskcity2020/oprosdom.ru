package repo

import (
	"context"

	"oprosdom.ru/msvc_auth/internal/models"
	repo_internal "oprosdom.ru/msvc_auth/internal/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	PhoneSend(ctx context.Context, v *models.ValidatedPhoneSendReq) error
}
