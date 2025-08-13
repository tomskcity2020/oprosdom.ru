package polls_repo

import (
	"context"

	polls_models "oprosdom.ru/core/internal/polls/models"
	repo_internal "oprosdom.ru/core/internal/polls/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	GetPolls(ctx context.Context) ([]*polls_models.Poll, error)
	Vote(ctx context.Context, m *polls_models.ValidVoteReq) error
	PollStats(ctx context.Context) ([]*polls_models.PollStats, error)
}
