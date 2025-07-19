package service

import (
	"oprosdom.ru/monolith/internal/dz/internal/biz"
	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	service_internal "oprosdom.ru/monolith/internal/dz/internal/service/internal"
)

type Service interface {
	MemberAdd(member *models.Member) error
	KvartiraAdd(kvartira *models.Kvartira) error
	MemberUpdate(member *models.Member) error
	KvartiraUpdate(kvartira *models.Kvartira) error
	MembersGet() ([]*models.Member, error)
	KvartirasGet() ([]*models.Kvartira, error)
	MemberGet(id string) (*models.Member, error)
	KvartiraGet(id string) (*models.Kvartira, error)
	RemoveById(mk string, id string) error
	CountData()
	RunParallel(modelsData []models.ModelInterface)
	RunSeq(modelsData []models.ModelInterface)
}

// фабрика будет вызывать другой конструктор из internal service
func NewServiceFactory() Service {
	repo, err := repo.GetRepoSingleton()
	if err != nil {
		panic(err)
	}
	biz := biz.NewBizFactory()
	return service_internal.NewCallInternalService(repo, biz)
}
