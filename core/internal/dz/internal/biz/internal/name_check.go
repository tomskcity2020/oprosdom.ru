package biz_internal

import (
	"errors"
	"strings"
)

func (b *BizStruct) nameCheck(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("empty name")
	}
	return nil
}
