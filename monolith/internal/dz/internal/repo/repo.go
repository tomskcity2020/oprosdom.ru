package repo

import (
	"oprosdom.ru/monolith/internal/dz/internal/models"
	repo_internal "oprosdom.ru/monolith/internal/dz/internal/repo/internal"
)

type RepositoryInterface interface {
	Save(m models.ModelInterface)
	LoadFromFile(fileName string)
	MembersInSliceNow() int
	KvartirasInSliceNow() int
	SaveToFile(fileName string)
	GetSliceMembers() []*models.Member
	GetSliceKvartiras() []*models.Kvartira
}

func NewRepoFactory() RepositoryInterface {
	return repo_internal.NewCallInternalRepo()
}
