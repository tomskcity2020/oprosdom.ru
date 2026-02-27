package models

import (
	"net"

	shared_validate "oprosdom.ru/shared/validate"
)

type UnsafePhoneSendReq struct {
	Phone     string `json:"phone"`
	UserAgent string `swaggerignore:"true"` // в phonesend используется для rate лимита
	Ip        string `swaggerignore:"true"` // в phonesend используется для rate лимита
}

func (u *UnsafePhoneSendReq) Validate() (*ValidatedPhoneSendReq, error) {

	var validatedPhoneSendReq ValidatedPhoneSendReq

	validPhone, _, err := shared_validate.PhoneValidate(u.Phone)
	if err != nil {
		return nil, err
	}

	validatedPhoneSendReq.Phone = validPhone
	// validatedPhoneSendReq.PhoneType = validPhoneType

	validatedPhoneSendReq.UserAgent = shared_validate.UserAgentSanitize(u.UserAgent)

	validIp, err := shared_validate.IpValidate(u.Ip)
	if err != nil {
		return nil, err
	}

	validatedPhoneSendReq.IP = validIp

	return &validatedPhoneSendReq, nil
}

type ValidatedPhoneSendReq struct {
	Phone string // в E.164
	//PhoneType string // mobile / landline / unknown
	UserAgent string // уже очищенный и обрезанный если len > 512
	IP        net.IP // проверенный net.IP (в нашем случае через nginx получаем x-real-ip)
}
