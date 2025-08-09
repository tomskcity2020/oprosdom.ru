package models

import (
	"errors"
	"net"

	shared_validate "oprosdom.ru/shared/validate"
)

type UnsafeCodeCheckReq struct {
	Phone     string `json:"phone"`
	Code      uint32 `json:"code"`
	UserAgent string // в codecheck используется для записи в postgresql
	Ip        string // в codecheck используется для записи в postgresql
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

	valid.UserAgent = shared_validate.UserAgentSanitize(u.UserAgent)

	validIp, err := shared_validate.IpValidate(u.Ip)
	if err != nil {
		return nil, err
	}

	valid.IP = validIp

	return &valid, nil
}

type ValidatedCodeCheckReq struct {
	Phone     string // в E.164
	Code      uint32
	UserAgent string // уже очищенный и обрезанный если len > 512
	IP        net.IP // проверенный net.IP (в нашем случае через nginx получаем x-real-ip)
}
