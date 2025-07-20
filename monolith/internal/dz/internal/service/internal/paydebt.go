package service_internal

import (
	"context"
	"errors"
	"fmt"
	"log"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (s *ServiceStruct) PayDebt(ctx context.Context, r *models.PayDebtRequest) (*models.PayDebtResponse, error) {

	// проводим первичную валидацию (когда структура заполнена включая id)
	// сначала отдельно проверяем id, так как из BasicValidation эта проверка исключена (ввиду генерации id на стороне базы данных)
	if err := s.biz.UuidCheck(r.MemberId); err != nil {
		return nil, fmt.Errorf("member id validation failed: %v", err.Error())
	}

	if err := s.biz.UuidCheck(r.KvartiraId); err != nil {
		return nil, fmt.Errorf("kvartira id validation failed: %v", err.Error())
	}

	if err := s.biz.DecimalCheck(r.Amount); err != nil {
		return nil, fmt.Errorf("amount validation failed: %v", err.Error())
	}

	log.Println(r)

	response, err := s.repo.PayDebt(ctx, r)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return response, nil

}
