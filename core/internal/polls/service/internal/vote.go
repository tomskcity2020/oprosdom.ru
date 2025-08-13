package service_internal

import (
	"context"

	polls_models "oprosdom.ru/core/internal/polls/models"
)

func (s *ServiceStruct) Vote(ctx context.Context, p *polls_models.ValidVoteReq) error {

	if err := s.repo.Vote(ctx, p); err != nil {
		return err
	}

	return nil

}
