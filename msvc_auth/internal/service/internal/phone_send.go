package service_internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/shared/models/pb"
)

func (s *ServiceStruct) PhoneSend(ctx context.Context, p *models.ValidatedPhoneSendReq) error {

	// 4) нет - делаем INCR key{code:value} и отправляем
	// mutex?

	// 1) проверяем антифлудом
	if err := s.antifloodPhone(ctx, p); err != nil {
		return err
	}

	// 2) смотрим есть ли ключ phone:+71231231234: если есть, значит получаем retry:+71231231234, но все равно проверяем если второй get вернет nil, то retry делаем 1 (такое поведение возможно когда между запросами истечет retry:+71231231234)
	// в отличие от создания записей - здесь атомарность нужна, потому что может возникнуть такая ситуация, что мы возьмем первый ключ, а второй истечет по ttl
	// в любом случае сначала чекаем на nil что пришло
	keys := []string{"phone:" + p.Phone, "retry:" + p.Phone} // внимание! ниже по ключу доступ, если меняем тут то с оглядкой
	values, err := s.ramRepo.GetFew(ctx, keys)
	if err != nil {
		return err
	}

	// 3) есть - значит есть и код внутри, значит не генерируем новый код, а отправляем тот же. При этом смотрим retry
	var existPhoneCode models.PhoneCode
	if values[0] != nil {
		strVal1, ok := values[0].(string)
		if !ok {
			err := "strVal1 is not a string"
			log.Println(err)
			return errors.New(err)
		}
		if err := json.Unmarshal([]byte(strVal1), &existPhoneCode); err != nil {
			log.Printf("json unmarshal failed: %v", err)
			return err
		}
	} else {
		fmt.Println("val1 отсутствует в Redis")
	}

	var retry uint32
	if values[1] != nil {
		strVal2, ok := values[1].(string)
		if !ok {
			err := "strVal2 is not a string"
			log.Println(err)
			return errors.New(err)
		}
		existRetryInt, err := strconv.Atoi(strVal2)
		if err != nil {
			log.Printf("atoi returns err: %v", err)
			return err
		}
		retry = uint32(existRetryInt)
	} else {
		// если по каким-то причинам retry не оказалось в redis, то приводим к 1 и логируем, так как такого события наступить не должно, чтоб если наступит - отреагировать
		log.Printf("retry records not exist in redis")
		retry = 1
	}

	// Логика в том, что повторные попытки скорее всего это неудачная доставка смс, > 3 это скорее всего уже баловство значит отправляем дешевым шлюзом. Также обращаем внимание на то, что отправляем тот же код преднамеренно: если будем каждую отправку менять код, то можем только запутать клиента: смски ненадежны и при более 2 попытках они могут доставиться клиенту не в том порядке в котором ожидается. Если мы даем окно 5 мин в пределах которого код будет один и тот же, то в этом нет ничего критичного. С учетом того, что на один тел 20 попыток ограничение в сутки и заменой кода через 5 минут - вполне безопасно.
	// если запись в редисе есть, это значит то, что одно сообщение уже точно улетело, а значит в этой итерации нужно ++
	// это алгоритм обеспечивает retry1,2,3,1,1,1,1,1...
	switch retry {
	case 1:
		retry = 2
	case 2:
		retry = 3
	default:
		retry = 1
	}

	msg := &pb.MsgCode{}
	msg.Phone = p.Phone

	if existPhoneCode.Code > 999 && existPhoneCode.Code < 10000 {

		msg.Code = existPhoneCode.Code
		msg.Retry = retry

		// увеличиваем count в редисе (не вставляем newretry через set иначе получим другую логику совершенно)
		if _, err := s.ramRepo.Incr(ctx, "retry:"+p.Phone); err != nil {
			log.Printf("retry incr failed: %v", err)
		}

	} else {

		// 3) генерируем код и создаем записи в redis'е (атомарность для записи неважна в этом конкретном случае, так как чтение будет происходить через достаточное время: пока отправится смс, пока доставится, пока клиент введет)
		code := uint32(rand.Intn(9000) + 1000)

		phoneCode := &models.PhoneCode{
			Phone: p.Phone,
			Code:  code,
		}

		jsonValue, err := json.Marshal(phoneCode)
		if err != nil {
			return err
		}

		msg.Code = code
		msg.Retry = 1

		// ttl 10 минут из расчета 3 кода по смс, 1 код по звонку + время на вводы между запросами = 10 минут
		if err := s.ramRepo.Set(ctx, "phone:"+p.Phone, jsonValue, 10*time.Minute); err != nil {
			return err // влияет на бизнес-логику, поэтому прерываем выполнение программы
		}

		// retry создаем тоже через set, чтоб удалялись записи заброшенные (по которым подтверждение кода не прошло - там удаляются phone и retry ключи)
		if err := s.ramRepo.Set(ctx, "retry:"+p.Phone, 1, 10*time.Minute); err != nil {
			log.Printf("cant set retry: %v", err) // не влияет на бизнес-логику, логируем и продолжаем
		}
		// if _, err := s.ramRepo.Incr(ctx, "retry:"+p.Phone); err != nil {
		// 	log.Printf("cant set retry: %v", err) // не влияет на бизнес-логику, логируем и продолжаем
		// }

	}

	if err := s.codeTransport.Send(ctx, msg); err != nil {
		log.Printf("Failed to send code: %v", err)
		return err
	}

	return nil

}
