package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type ModelInterface interface {
	Type() string
}

// func NewUserFactory(name string, phone string, community int) ModelInterface {
// 	return NewMember(name, phone, community)
// }

// func NewKvartiraFactory(number string, komnat int) ModelInterface {
// 	return NewKvartira(number, komnat)
// }

type Member struct {
	Id        string    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Community int    `json:"community"`
}

func (m *Member) BasicValidate() error {
	if _, err := uuid.Parse(m.Id); err !=nil {
		return errors.New("incorrect id")
	}

	m.Name = strings.TrimSpace(m.Name)
	if m.Name == "" {
		return errors.New("empty name")
	}

	m.Phone = strings.TrimSpace(m.Phone)
	if m.Phone == "" {
		return errors.New("empty phone")
	}

	if m.Community <= 0 {
		return errors.New("incorrect Community")
	}

	return nil
}

func (m *Member) Type() string {
	return "member"
}

func NewMember(name string, phone string, community int) *Member {
	return &Member{
		Name:      name,
		Phone:     phone,
		Community: community,
	}
}

type Kvartira struct {
	Number string
	Komnat int
}

func (obj *Kvartira) Type() string {
	return "kvartira"
}

func NewKvartira(number string, komnat int) *Kvartira {
	return &Kvartira{
		Number: number,
		Komnat: komnat,
	}
}
