package repo

import (
	"sync"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	repo_internal "oprosdom.ru/monolith/internal/dz/internal/repo/internal"
)

// глобальная переменная
var (
	globalRepo  RepositoryInterface
	repoOnce    sync.Once
	repoInitErr error
)

// вместо NewRepoFactory реализовываем тн синглтон (когда структура создается только 1 раз глобально - нам это нужно чтоб между запросами post/delete запоминать инфу в слайсах)
func GetRepoSingleton() (RepositoryInterface, error) {
	repoOnce.Do(func() {
		globalRepo, repoInitErr = repo_internal.NewCallInternalRepo()
	})
	return globalRepo, repoInitErr
}

type RepositoryInterface interface {
	Save(m models.ModelInterface) error
	LoadFromFile(fileName string)
	UpdateFile(m models.ModelInterface) error
	UpdateSlice(m models.ModelInterface) error
	MembersInSliceNow() int
	KvartirasInSliceNow() int
	SaveToFile(m models.ModelInterface) error
	GetSliceMembers() []*models.Member
	GetSliceKvartiras() []*models.Kvartira
}

// func NewRepoFactory() RepositoryInterface {
// 	return repo_internal.NewCallInternalRepo()
// }
