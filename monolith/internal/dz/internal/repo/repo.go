package repo

import (
	"context"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	repo_internal "oprosdom.ru/monolith/internal/dz/internal/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	KvartiraAdd(ctx context.Context, k *models.Kvartira) error
	KvartiraGetById(ctx context.Context, id string) (*models.Kvartira, error)
	KvartiraUpdate(ctx context.Context, k *models.Kvartira) error
	KvartirasGet(ctx context.Context) ([]*models.Kvartira, error)
	MemberAdd(ctx context.Context, m *models.Member) error
	MemberGetById(ctx context.Context, id string) (*models.Member, error)
	MemberUpdate(ctx context.Context, m *models.Member) error
	MembersGet(ctx context.Context) ([]*models.Member, error)
	DeleteById(ctx context.Context, id string, mk string) error
	PayDebt(ctx context.Context, r *models.PayDebtRequest) (*models.PayDebtResponse, error)
}
