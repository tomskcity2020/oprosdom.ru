package service_internal

import (
	"context"
	"errors"
	"hash/crc32"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"oprosdom.ru/msvc_auth/internal/models"
)

func (s *ServiceStruct) antifloodPhone(ctx context.Context, p *models.ValidatedPhoneSendReq) error {

	// правила:
	// 1) один и тот же номер может не более 10 раз указывать за 24 часа (ttl и по count смотреть) (чтоб никого не заспамили критически через нас)
	// 2) клиент с одинаковым useragent и ip имеет право делать запросы не чаще 1 раза в 20 сек (ttl) (защита от тех кто не умеет useragent менять, но для нас это хорошее уточнее для антифлуда)
	// 3) с одного и того же ip можно делать запрос не чаще 1 раза в 5 секунд (ttl) (защита от продуманых злодеев, но грубая на первое время)

	// TODO нужно усложнять правила, делать умнее, не такими грубыми
	// + добавить капчу на подозрительные запросы

	// хэш делаем через CRC32 он самый быстрый, но вероятность коллизий выше. Но в нашем случае ключ состоит из хэш ua это только половина ключа, вторая половина - ip адрес, который почти полностью нивелирует опасность от коллизий
	uaHash := crc32.ChecksumIEEE([]byte(p.UserAgent))
	uaStr := strconv.FormatUint(uint64(uaHash), 10)

	rate_phone := "rate_phone:" + p.Phone
	rate_uaip := "rate_uaip:" + uaStr + ":" + p.IP.String()
	rate_ip := "rate_ip:" + p.IP.String()

	countStr, err := s.ramRepo.Get(ctx, rate_phone)
	if errors.Is(err, redis.Nil) {
		// ключ не найден, записываем первое значение и продолжаем выполнение
		s.ramRepo.Set(ctx, rate_phone, 1, 24*time.Hour)
	} else if err != nil {
		// какая-то другая ошибка редиса, логируем
		log.Printf("redis error occured: %v", err)
		return err
	} else {
		// ключ найден, проверяем не превышено ли кол-во count
		count, err := strconv.Atoi(countStr)
		if err != nil {
			log.Printf("atoi failed on antiflood_phone: %v", err)
			return err
		}
		// запись в редисе с ttl 24 часа, если больше 20 раз номер отправляли - ошибка, если меньше - делаем incr и продолжаем выполнение
		if count > 20 {
			return errors.New("count exceeded 24h limit")
		} else {
			s.ramRepo.Incr(ctx, rate_phone)
		}
	}

	_, err = s.ramRepo.Get(ctx, rate_uaip)
	if errors.Is(err, redis.Nil) {
		// ключ не найден, записываем первое значение и продолжаем выполнение
		s.ramRepo.Set(ctx, rate_uaip, 1, 20*time.Second)
	} else if err != nil {
		// какая-то другая ошибка редиса, логируем
		log.Printf("redis error occured: %v", err)
		return err
	} else {
		// ключ найден - ошибка
		return errors.New("rate_uaip exceeded")
	}

	_, err = s.ramRepo.Get(ctx, rate_ip)
	if errors.Is(err, redis.Nil) {
		// ключ не найден, записываем первое значение и продолжаем выполнение
		s.ramRepo.Set(ctx, rate_ip, 1, 5*time.Second)
	} else if err != nil {
		// какая-то другая ошибка редиса, логируем
		log.Printf("redis error occured: %v", err)
		return err
	} else {
		// ключ найден - ошибка
		return errors.New("rate_ip exceeded")
	}

	return nil

}
