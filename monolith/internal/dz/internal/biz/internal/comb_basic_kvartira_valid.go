package biz_internal

import (
	"errors"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (b *BizStruct) BasicKvartiraValidation(kvartira *models.Kvartira) error {
	// if err := b.UuidCheck(kvartira.Id); err != nil {
	// 	return errors.New(err.Error())
	// }
	if err := b.kvNumberCheck(kvartira.Number); err != nil {
		return errors.New(err.Error())
	}
	if err := b.kvKomnatCheck(kvartira.Komnat); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
