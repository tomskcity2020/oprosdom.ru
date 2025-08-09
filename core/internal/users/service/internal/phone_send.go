package users_service_internal

import (
	"context"
	"errors"

	users_models "oprosdom.ru/core/internal/users/models"
)

func (s *ServiceStruct) PhoneSend(ctx context.Context, p *users_models.ValidatedPhoneSendReq) error {

	// TODO
	// таблица phonesend в postgresql создана для длительного хранения телефонов (аналитка + будущий функционал)
	// первой линией обороны нужно сделать redis: перед каждым запросом чекаем записи, если по своду антифлуд-правил все норм, то пишем в postgresql
	// если не пройдет антифлуд, то отдаем команду что типа все ок, код отправлен

	if err := s.repo.PhoneSend(ctx, p); err != nil {
		return errors.New(err.Error())
	}

	return nil

}
