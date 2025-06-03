package repo

import (
	"oprosdom.ru/monolith/internal/dz/internal/models"
	repo_internal "oprosdom.ru/monolith/internal/dz/internal/repo/internal"
)

type RepositoryInterface interface {
	Save(m models.ModelInterface)
	Check()
	Show(t string) int
}

func NewRepoFactory() RepositoryInterface {
	return repo_internal.NewCallInternalRepo()
}
