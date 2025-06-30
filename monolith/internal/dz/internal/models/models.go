package models

type ModelInterface interface {
	Type() string
}

func NewUserFactory(name string, phone string, community int) ModelInterface {
	return NewMember(name, phone, community)
}

func NewKvartiraFactory(number string, komnat int) ModelInterface {
	return NewKvartira(number, komnat)
}

type Member struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Community int    `json:"community"`
}

func (obj *Member) Type() string {
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
