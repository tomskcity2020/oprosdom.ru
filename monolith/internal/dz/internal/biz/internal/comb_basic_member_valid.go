package biz_internal

import (
	"errors"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func (b *BizStruct) BasicMemberValidation(member *models.Member) error {
	// if err := b.UuidCheck(member.Id); err != nil {
	// 	return errors.New(err.Error())
	// }
	if err := b.nameCheck(member.Name); err != nil {
		return errors.New(err.Error())
	}
	if err := b.phoneCheck(member.Phone); err != nil {
		return errors.New(err.Error())
	}
	if err := b.communityIdCheck(member.Community); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
