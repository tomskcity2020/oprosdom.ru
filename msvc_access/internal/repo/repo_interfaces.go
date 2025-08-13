package repo

import (
	"context"
	"time"

	repo_internal "oprosdom.ru/msvc_access/internal/repo/internal"
)

func NewRamRepoFactory(ctx context.Context, addr string) (RamRepoInterface, error) {
	return repo_internal.NewRedis(ctx, addr)
}

type RamRepoInterface interface {
	Close()
	Set(ctx context.Context, k string, v any, ttl time.Duration) error
}
