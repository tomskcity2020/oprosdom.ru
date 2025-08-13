package polls_service

import (
	"context"

	polls_biz "oprosdom.ru/core/internal/polls/biz"
	polls_models "oprosdom.ru/core/internal/polls/models"
	polls_repo "oprosdom.ru/core/internal/polls/repo"
	service_internal "oprosdom.ru/core/internal/polls/service/internal"
)

type PollsService interface {
	GetPolls(ctx context.Context) ([]*polls_models.Poll, error)
	Vote(ctx context.Context, p *polls_models.ValidVoteReq) error
	PollStats(ctx context.Context) ([]*polls_models.PollStats, error)
}

// обертка чтоб вызывать из internal и возвращать интерфейс
func NewServiceFactory(repo polls_repo.RepositoryInterface) PollsService {
	biz := polls_biz.NewBizFactory()
	return service_internal.NewCallInternalService(repo, biz)
}
