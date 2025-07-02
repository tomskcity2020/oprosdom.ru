package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type ModelInterface interface {
}

// func NewUserFactory(name string, phone string, community int) ModelInterface {
// 	return NewMember(name, phone, community)
// }

// func NewKvartiraFactory(number string, komnat int) ModelInterface {
// 	return NewKvartira(number, komnat)
// }

type Member struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Community int    `json:"community"`
}

// метод CreateUuid будет использоваться только с POST при создании нового member. А BasicValidate тогда можем и для POST и для PUT использовать
func (m *Member) CreateUuid() error {
	m.Id = uuid.NewString()
	return nil
}

func (m *Member) BasicValidate() error {
	if _, err := uuid.Parse(m.Id); err != nil {
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

// func NewMember(name string, phone string, community int) *Member {
// 	return &Member{
// 		Name:      name,
// 		Phone:     phone,
// 		Community: community,
// 	}
// }

type Kvartira struct {
	Id     string `json:"id"`
	Number string `json:"number"`
	Komnat int    `json:"komnat"`
}

// func NewKvartira(number string, komnat int) *Kvartira {
// 	return &Kvartira{
// 		Number: number,
// 		Komnat: komnat,
// 	}
// }

func (m *Kvartira) CreateUuid() error {
	m.Id = uuid.NewString()
	return nil
}

func (m *Kvartira) BasicValidate() error {
	if _, err := uuid.Parse(m.Id); err != nil {
		return errors.New("incorrect id")
	}

	m.Number = strings.TrimSpace(m.Number)
	if m.Number == "" {
		return errors.New("empty number")
	}

	if m.Komnat <= 0 {
		return errors.New("incorrect komnat")
	}

	return nil
}
