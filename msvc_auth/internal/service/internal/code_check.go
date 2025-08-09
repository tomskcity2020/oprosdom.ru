package service_internal

import (
	"context"
	"errors"
	"log"
	"time"

	"oprosdom.ru/msvc_auth/internal/models"
)

func (s *ServiceStruct) CodeCheck(ctx context.Context, p *models.ValidatedCodeCheckReq) error {

	// одним запросом получаем 3 записи из редиса:
	// phone:+71231231234 (value: phone, code)
	// code_daily_count:+71231231234
	// code_attempt:+71231231234

	keys := []string{"phone:" + p.Phone, "code_daily_count:" + p.Phone, "code_attempt:" + p.Phone}
	values, err := s.ramRepo.GetFew(ctx, keys)
	if err != nil {
		return err
	}

	// если не будет всех ключей, то вернется слайс из 3 nil
	// values[0]: phone:+71231231234 (value: phone, code)

	// сразу отдаем ошибку если не найдена запись телефон/код
	if values[0] == nil {
		return errors.New("record_not_exists")
	}

	existPhoneCode, err := s.parsePhoneCode(values[0])
	if err != nil {
		log.Printf("failed to parse phone code: %v", err)
		return err
	}

	// values[1]: code_daily_count:+71231231234 count
	if values[1] != nil {
		codeDailycount, err := s.parseUint32(values[1])
		if err != nil {
			return err
		}

		// проверяем чтоб было не больше 20 попыток ввода кода по этому номеру (исходим из того, что человек может на нескольких устройствах аутентифицироваться + на ошибки возможные запас)
		if codeDailycount > 20 {
			return errors.New("daily_limit_exeeded")
		}

		// если ключ code_daily_count есть в редисе то увеличиваем значение через INCR, через SET нельзя иначе ttl установится заново
		if _, err := s.ramRepo.Incr(ctx, "code_daily_count:"+p.Phone); err != nil {
			return err // влияет на бизнес-логику, поэтому прерываем выполнение программы
		}

	} else {
		// если ключа code_daily_count нет в редисе, то создаем через SET (устанавливаем ttl)
		if err := s.ramRepo.Set(ctx, "code_daily_count:"+p.Phone, 0, 24*time.Hour); err != nil {
			return err // влияет на бизнес-логику, поэтому прерываем выполнение программы
		}
	}

	// values[2]: code_attempt:+71231231234 count
	// code_attempt:+71231231234 устанавливаем в phone_send в одно время с phone:+71231231234 для того, чтобы конкретно на 1 код давать не больше 10 попыток. Потому что если рассинхрон произойдет, то непредсказумая логика будет. Может ничего негативного не будет, но для порядка делаем так.
	if values[2] != nil {
		codeAttempt, err := s.parseUint32(values[2])
		if err != nil {
			return err
		}

		// разрешаем только 10 попыток для ввода кода (активный код имеет такой же ttl как и code_attempt так как созданы в одно время в phone_send)
		if codeAttempt > 9 {
			return errors.New("try_again_10min")
		}

		// если ключ code_attempt есть в редисе то увеличиваем значение через INCR, через SET нельзя иначе ttl установится заново
		if _, err := s.ramRepo.Incr(ctx, "code_attempt:"+p.Phone); err != nil {
			return err // влияет на бизнес-логику, поэтому прерываем выполнение программы
		}
	}
	// если ключа code_attempt нет в редисе, то создавать здесь через SET не нужно, потому что code_attempt создается в phone_send вместе с phone:phone для синхронного ttl

	if p.Code != existPhoneCode.Code {
		return errors.New("code_not_accepted")
	}

	log.Println("CONFIRMED")

	return nil

}
