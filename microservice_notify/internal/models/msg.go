package models

import (
	"errors"
	"strings"

	shared_validate "oprosdom.ru/shared/validate"
)

type UnsafeMsg struct {
	Urgent      bool
	Type        string
	Phone       string `json:"phone"`
	MessageText string
	Retry       int32
}

func (u *UnsafeMsg) Validate() (*ValidatedMsg, error) {

	var validatedMsg ValidatedMsg

	validatedMsg.Urgent = u.Urgent

	validPhone, validPhoneType, err := shared_validate.PhoneValidate(u.Phone)
	if err != nil {
		return nil, err
	}

	validatedMsg.Type = validPhoneType
	validatedMsg.Phone = validPhone

	if strings.TrimSpace(u.MessageText) == "" {
		return nil, errors.New("message text cannot be empty")
	}
	validatedMsg.MessageText = u.MessageText

	if u.Retry <= 0 || u.Retry > 3 {
		return nil, errors.New("retry must be between 1 and 3")
	}
	validatedMsg.Retry = u.Retry

	return &validatedMsg, nil
}

type ValidatedMsg struct {
	Urgent      bool
	Type        string
	Phone       string `json:"phone"`
	MessageText string
	Retry       int32
}
