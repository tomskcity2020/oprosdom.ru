package repo

import (
	"oprosdom.ru/monolith/internal/dz6/internal/models"
	repo_internal "oprosdom.ru/monolith/internal/dz6/internal/repo/internal"
)

type RepositoryInterface interface {
	Save(m models.ModelInterface)
	Show(t string) int
}

func NewRepoFactory() RepositoryInterface {
	return repo_internal.NewCallInternalRepo()
}
