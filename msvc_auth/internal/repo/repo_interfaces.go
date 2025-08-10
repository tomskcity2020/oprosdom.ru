package repo

import (
	"context"
	"time"

	"oprosdom.ru/msvc_auth/internal/models"
	repo_internal "oprosdom.ru/msvc_auth/internal/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	AddSignedToken(ctx context.Context, v *models.ValidatedCodeCheckReq, k *models.KeyData) error
}

func NewRamRepoFactory(ctx context.Context, addr string) (RamRepoInterface, error) {
	return repo_internal.NewRedis(ctx, addr)
}

type RamRepoInterface interface {
	Close()
	Incr(ctx context.Context, k string) (int64, error)
	Set(ctx context.Context, k string, v any, ttl time.Duration) error
	Get(ctx context.Context, k string) (string, error)
	GetFew(ctx context.Context, keys []string) ([]any, error)
	Del(ctx context.Context, k string) (int64, error)
}
