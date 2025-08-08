package shared_validate

import (
	"errors"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func PhoneValidate(phone string) (formattedPhone string, phoneType string, err error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "", "", errors.New("empty_phone")
	}

	tel, err := phonenumbers.Parse(phone, "RU")
	if err != nil {
		return "", "", errors.New("incorrect_format_phone")
	}

	if !phonenumbers.IsValidNumberForRegion(tel, "RU") {
		return "", "", errors.New("not_valid_ru_phone_number")
	}

	switch phonenumbers.GetNumberType(tel) {
	case phonenumbers.MOBILE:
		phoneType = "mobile"
	case phonenumbers.FIXED_LINE, phonenumbers.FIXED_LINE_OR_MOBILE:
		phoneType = "landline"
	default:
		phoneType = "unknown"
	}

	formattedPhone = phonenumbers.Format(tel, phonenumbers.E164)

	return formattedPhone, phoneType, nil
}
