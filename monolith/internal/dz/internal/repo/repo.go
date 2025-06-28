package repo

import (
	"context"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	repo_internal "oprosdom.ru/monolith/internal/dz/internal/repo/internal"
)

type RepositoryInterface interface {
	Save(saveCtx context.Context, m models.ModelInterface)
	GetSliceMembers() []*models.Member
	GetSliceKvartiras() []*models.Kvartira
}

func NewRepoFactory() RepositoryInterface {
	return repo_internal.NewCallInternalRepo()
}
