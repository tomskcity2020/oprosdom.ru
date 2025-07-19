package biz

import (
	biz_internal "oprosdom.ru/monolith/internal/dz/internal/biz/internal"
	"oprosdom.ru/monolith/internal/dz/internal/models"
)

type BizInterface interface {
	UuidCreate() (string, error)
	UuidCheck(id string) error
	BasicMemberValidation(member *models.Member) error
	BasicKvartiraValidation(kvartira *models.Kvartira) error
}

func NewBizFactory() BizInterface {
	return biz_internal.NewCallInternalBiz()
}
