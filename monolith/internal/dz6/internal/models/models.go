package models

import (
	models_internal "oprosdom.ru/monolith/internal/dz6/internal/models/internal"
)

type ModelInterface interface {
	Type() string
}

func NewUserFactory(name string, phone string, community int) ModelInterface {
	return models_internal.NewMember(name, phone, community)
}

func NewKvartiraFactory(number string, komnat int) ModelInterface {
	return models_internal.NewKvartira(number, komnat)
}
