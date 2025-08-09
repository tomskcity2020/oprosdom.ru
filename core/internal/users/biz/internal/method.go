package users_biz_internal

// func (b *BizStruct) PhoneNumberCheck(phone string) error {
// 	phone = strings.TrimSpace(phone)
// 	if phone == "" {
// 		return errors.New("empty_phone")
// 	}

// 	tel, err := phonenumbers.Parse(phone, "RU")
// 	if err != nil {
// 		return errors.New("incorrect_format_phone")
// 	}

// 	if !phonenumbers.IsValidNumber(tel) {
// 		return errors.New("not_valid__ru_phone_number")
// 	}

// 	return nil
// }
