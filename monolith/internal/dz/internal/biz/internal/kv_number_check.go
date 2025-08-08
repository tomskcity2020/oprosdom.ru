package biz_internal

import (
	"errors"
	"strings"
)

func (b *BizStruct) kvNumberCheck(number string) error {
	number = strings.TrimSpace(number)
	if number == "" {
		return errors.New("empty kv number")
	}
	return nil
}
