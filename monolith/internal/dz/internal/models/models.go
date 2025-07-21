package models

type ModelInterface interface{}

type Member struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Community int    `json:"community"`
}

type Kvartira struct {
	Id     string `json:"id"`
	Number string `json:"number"`
	Komnat int    `json:"komnat"`
}
