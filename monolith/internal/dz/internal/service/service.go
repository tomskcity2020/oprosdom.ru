package service

import (
	"context"

	"oprosdom.ru/monolith/internal/dz/internal/biz"
	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	service_internal "oprosdom.ru/monolith/internal/dz/internal/service/internal"
)

type Service interface {
	KvartiraAdd(ctx context.Context, kvartira *models.Kvartira) error
	KvartiraGet(ctx context.Context, id string) (*models.Kvartira, error)
	KvartiraUpdate(ctx context.Context, k *models.Kvartira) error
	KvartirasGet(ctx context.Context) ([]*models.Kvartira, error)
	MemberAdd(ctx context.Context, m *models.Member) error
	MemberGet(ctx context.Context, id string) (*models.Member, error)
	MemberUpdate(ctx context.Context, m *models.Member) error
	MembersGet(ctx context.Context) ([]*models.Member, error)
	RemoveById(ctx context.Context, id string, mk string) error
	PayDebt(ctx context.Context, r *models.PayDebtRequest) (*models.PayDebtResponse, error)
	// CountData()
	// RunParallel(modelsData []models.ModelInterface)
	// RunSeq(modelsData []models.ModelInterface)
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory(repo repo.RepositoryInterface) Service {
	// repo, err := repo.GetRepoSingleton()
	// if err != nil {
	// 	panic(err)
	// }

	biz := biz.NewBizFactory()
	return service_internal.NewCallInternalService(repo, biz)
}
