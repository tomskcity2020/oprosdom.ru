package biz_internal

import (
	"errors"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func (b *BizStruct) phoneCheck(phone string) error {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return errors.New("empty phone")
	}

	mobile, err := phonenumbers.Parse(phone, "RU")
	if err != nil {
		return errors.New("incorrect format phone")
	}

	if !phonenumbers.IsValidNumber(mobile) {
		return errors.New("not valid phone number")
	}

	return nil
}
