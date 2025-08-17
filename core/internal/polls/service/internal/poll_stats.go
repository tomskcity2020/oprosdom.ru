package service_internal

import (
	"context"

	polls_models "oprosdom.ru/core/internal/polls/models"
)

func (s *ServiceStruct) PollStats(ctx context.Context) ([]*polls_models.PollStats, error) {

	m, err := s.repo.PollStats(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil

}
