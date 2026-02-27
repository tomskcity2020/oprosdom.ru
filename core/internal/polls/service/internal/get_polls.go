package service_internal

import (
	"context"

	polls_models "oprosdom.ru/core/internal/polls/models"
)

func (s *ServiceStruct) GetPolls(ctx context.Context) ([]*polls_models.Poll, error) {

	m, err := s.repo.GetPolls(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil

}
