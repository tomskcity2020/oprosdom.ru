package users_biz

import (
	users_biz_internal "oprosdom.ru/monolith/internal/users/biz/internal"
	users_models "oprosdom.ru/monolith/internal/users/models"
)

type BizInterface interface {
	UuidCreate() (string, error)
	UuidCheck(id string) error
	BasicMemberValidation(member *users_models.Member) error
	BasicKvartiraValidation(kvartira *models.Kvartira) error
	DecimalCheck(id string) error
}

func NewBizFactory() BizInterface {
	return users_biz_internal.NewCallInternalBiz()
}
