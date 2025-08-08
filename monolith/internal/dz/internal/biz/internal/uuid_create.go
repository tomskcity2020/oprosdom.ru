package biz_internal

import (
	"errors"

	"github.com/google/uuid"
)

// сразу NewString не делаем, так как не возвращает ошибку в отличие от UUID.
func (b *BizStruct) UuidCreate() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.New("uuid generation failed")
	}
	return id.String(), nil
}
