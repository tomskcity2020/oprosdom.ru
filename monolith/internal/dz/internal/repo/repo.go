package repo

import (
	"sync"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	repo_internal "oprosdom.ru/monolith/internal/dz/internal/repo/internal"
)

// глобальная переменная
var (
	globalRepo RepositoryInterface
	repoOnce   sync.Once
)

// вместо NewRepoFactory реализовываем тн синглтон (когда структура создается только 1 раз глобально - нам это нужно чтоб между запросами post/delete запоминать инфу в слайсах)
func GetRepoSingleton() RepositoryInterface {
	repoOnce.Do(func() {
		globalRepo = repo_internal.NewCallInternalRepo()
	})
	return globalRepo
}

type RepositoryInterface interface {
	Save(m models.ModelInterface) error
	LoadFromFile(fileName string)
	MembersInSliceNow() int
	KvartirasInSliceNow() int
	SaveToFile(fileName string)
	GetSliceMembers() []*models.Member
	GetSliceKvartiras() []*models.Kvartira
}

// func NewRepoFactory() RepositoryInterface {
// 	return repo_internal.NewCallInternalRepo()
// }
