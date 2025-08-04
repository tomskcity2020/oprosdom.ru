package service_internal

import (
	"context"
	"errors"
	"log"

	"oprosdom.ru/microservice_auth/internal/models"
	"oprosdom.ru/shared/models/pb"
)

func (s *ServiceStruct) PhoneSend(ctx context.Context, p *models.ValidatedPhoneSendReq) error {

	// TODO Redis antiflood records
	// таблица phonesend в postgresql создана для длительного хранения телефонов (аналитка + будущий функционал)
	// первой линией обороны нужно сделать redis: перед каждым запросом чекаем записи, если по своду антифлуд-правил все норм, то пишем в postgresql
	// если не пройдет антифлуд, то отдаем команду что типа все ок, код отправлен - чтоб сбивать злоумышленника с толку

	if err := s.repo.PhoneSend(ctx, p); err != nil {
		return errors.New(err.Error())
	}

	msg := &pb.MsgCode{
		Urgent:      true,
		Type:        "sms",
		PhoneNumber: "+79994951548",
		Message:     "OprosDom.ru ваш код 1234",
		Retry:       1,
	}

	if err := s.codeTransport.Send(ctx, msg); err != nil {
		log.Printf("Failed to send code: %v", err)
	} else {
		log.Println("Code sent successfully")
	}

	return nil

}
