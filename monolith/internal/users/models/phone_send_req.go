package users_models

import (
	"net"

	shared_validate "oprosdom.ru/shared/validate"
)

type UnsafePhoneSendReq struct {
	Phone     string `json:"phone"`
	UserAgent string
	Ip        string
}

func (u *UnsafePhoneSendReq) Validate() (*ValidatedPhoneSendReq, error) {

	var validatedPhoneSendReq ValidatedPhoneSendReq

	validPhone, validPhoneType, err := shared_validate.PhoneValidate(u.Phone)
	if err != nil {
		return nil, err
	}

	validatedPhoneSendReq.Phone = validPhone
	validatedPhoneSendReq.PhoneType = validPhoneType

	validatedPhoneSendReq.UserAgent = shared_validate.UserAgentSanitize(u.UserAgent)

	validIp, err := shared_validate.IpValidate(u.Ip)
	if err != nil {
		return nil, err
	}

	validatedPhoneSendReq.IP = validIp

	return &validatedPhoneSendReq, nil
}

type ValidatedPhoneSendReq struct {
	Phone     string // в E.164
	PhoneType string // mobile / landline / unknown
	UserAgent string // уже очищенный и обрезанный если len > 512
	IP        net.IP // проверенный net.IP (в нашем случае через nginx получаем x-real-ip)
}
