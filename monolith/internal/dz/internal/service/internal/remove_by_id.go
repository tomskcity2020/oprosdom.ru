package service_internal

import (
	"errors"
)

func (s *ServiceStruct) RemoveById(mk string, id string) error {

	var filename string

	switch mk {
	case "member":
		filename = "members"
	case "kvartira":
		filename = "kvartiras"
	default:
		return errors.New("неправильный запрос")
	}

	// чекаем что id является корректным uuid
	if err := s.biz.UuidCheck(id); err != nil {
		return errors.New("неправильный id")
	}

	if err := s.repo.RemoveFromFile(filename, id); err != nil {
		return errors.New("не удалось удалить из файла")
	}

	switch filename {
	case "members":
		if err := s.repo.RemoveMemberSlice(id); err != nil {
			return errors.New("не удалось удалить из слайса жителей")
		}
	case "kvartiras":
		if err := s.repo.RemoveKvartiraSlice(id); err != nil {
			return errors.New("не удалось удалить из слайса квартир")
		}
	default:
		return errors.New("внутренняя ошибка")
	}

	return nil

}
