package users_service_internal

import (
	"context"
	"errors"
	"fmt"

	"oprosdom.ru/monolith/internal/users/models"
)

func (s *ServiceStruct) PhoneSend(ctx context.Context, p *users_models) error {

	// сначала проверяем телефон, потому что если он невалиден, то смысла тратить время на проверку остального нет
	if err := s.biz.PhoneNumberCheck(p.Phone); err != nil {
		return err
	}
	
	




	if err := s.biz.UuidCheck(k.Id); err != nil {
		return fmt.Errorf("id validation failed: %v", err.Error())
	}

	if err := s.biz.BasicKvartiraValidation(k); err != nil {
		return fmt.Errorf("basic validation failed: %v", err.Error())
	}

	if err := s.repo.KvartiraUpdate(ctx, k); err != nil {
		return errors.New(err.Error())
	}

	return nil

}
