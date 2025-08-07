package service_internal

import (
	"context"

	"oprosdom.ru/msvc_auth/internal/models"
)

func (s *ServiceStruct) PhoneSend(ctx context.Context, p *models.ValidatedPhoneSendReq) error {

	// ключ: телефон
	// данные: код
	// алгоритм:
	// 1) проверяем антифлудом
	// 2) если все ок, то далее смотрим есть ли ключ:
	// 3) есть - значит есть и код внутри, значит не генерируем новый код, а отправляем тот же. При этом смотрим count. Если count до 3 вкл , то == retry. Count от 4 до 5 вкл == retry3, count > 5 == retry1. Логика в том, что повторные попытки скорее всего это неудачная доставка смс, 4 и 5 на запас, а свыше 5 это скорее всего уже баловство значит отправляем дешевым шлюзом. Также обращаем внимание на то, что отправляем тот же код преднамеренно: если будем каждую отправку менять код, то можем только запутать клиента: смски ненадежны и при более 2 попытках они могут доставиться клиенту не в том порядке в котором ожидается. Если мы даем окно 5 мин в пределах которого код будет один и тот же, то в этом нет ничего критичного. С учетом того, что на один тел 20 попыток ограничение в сутки и заменой кода через 5 минут - вполне безопасно.
	// 4) нет - делаем INCR key{code:value} и отправляем
	// 5) внезависимости от того есть ключ или нет, после отправки добавляем ключ user_agent:ip ttl 15s (всего делать самый наидешевейший хэш ua) и ip 5s ttl.
	// mutex?

	// 1) проверяем антифлудом
	if err := s.antifloodPhone(ctx, p); err != nil {
		return err
	}

	// if err := s.repo.PhoneSend(ctx, p); err != nil {
	// 	return errors.New(err.Error())
	// }

	// msg := &pb.MsgCode{
	// 	Phone: "+73822724299",
	// 	Code:  1234,
	// 	Retry: 3,
	// }

	// if err := s.codeTransport.Send(ctx, msg); err != nil {
	// 	log.Printf("Failed to send code: %v", err)
	// } else {
	// 	log.Println("Code sent successfully")
	// }

	return nil

}
