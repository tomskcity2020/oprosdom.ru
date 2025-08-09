package biz_internal

import (
	"errors"

	"github.com/google/uuid"
)

func (b *BizStruct) UuidCheck(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("incorrect id")
	}
	return nil
}
