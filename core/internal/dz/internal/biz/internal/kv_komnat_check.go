package biz_internal

import "errors"

func (b *BizStruct) kvKomnatCheck(komnat int) error {
	if komnat <= 0 {
		return errors.New("not valid quantity komnat")
	}
	return nil
}