package models_internal

type Member struct {
	Name      string
	Phone     string
	Community int
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
