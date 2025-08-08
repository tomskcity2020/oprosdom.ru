package models

import (
	"errors"

	shared_validate "oprosdom.ru/shared/validate"
)

type UnsafeMsg struct {
	Phone string
	Code  uint32
	Retry uint32
}

func (u *UnsafeMsg) Validate() (*ValidatedMsg, error) {

	var validatedMsg ValidatedMsg

	validPhone, validPhoneType, err := shared_validate.PhoneValidate(u.Phone)
	if err != nil {
		return nil, err
	}

	validatedMsg.Type = validPhoneType
	validatedMsg.Phone = validPhone

	if u.Code < 1000 || u.Code > 9999 {
		return nil, errors.New("code must be 1000 to 9999")
	}
	validatedMsg.Code = u.Code

	if u.Retry < 1 || u.Retry > 3 {
		return nil, errors.New("retry must be 1 or 2 or 3")
	}
	validatedMsg.Retry = u.Retry

	return &validatedMsg, nil
}

type ValidatedMsg struct {
	Type  string
	Phone string `json:"phone_number"`
	Code  uint32
	Retry uint32
}
