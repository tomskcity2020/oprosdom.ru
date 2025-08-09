package models

import (
	"errors"

	shared_validate "oprosdom.ru/shared/validate"
)

type UnsafeCodeCheckReq struct {
	Phone string `json:"phone"`
	Code  uint32 `json:"code"`
}

func (u *UnsafeCodeCheckReq) Validate() (*ValidatedCodeCheckReq, error) {

	var valid ValidatedCodeCheckReq

	validPhone, _, err := shared_validate.PhoneValidate(u.Phone)
	if err != nil {
		return nil, err
	}

	valid.Phone = validPhone

	if u.Code < 1000 || u.Code > 9999 {
		return nil, errors.New("code not valid")
	}

	valid.Code = u.Code

	return &valid, nil
}

type ValidatedCodeCheckReq struct {
	Phone string // в E.164
	Code  uint32
}
