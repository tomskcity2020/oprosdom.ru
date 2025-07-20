package models

type ModelInterface interface{}

type Member struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Community int    `json:"community"`
	Balance   string `json:"balance"` // все вычисления будут в базе данных, поэтому string
}

type Kvartira struct {
	Id     string `json:"id"`
	Number string `json:"number"`
	Komnat int    `json:"komnat"`
	Debt   string `json:"debt"`
}

type PayDebtRequest struct {
	MemberId   string `json:"id"`
	KvartiraId string `json:"kvartira_id"`
	Amount     string `json:"amount"`
}

type PayDebtResponse struct {
	NewBalance string `json:"new_balance"`
	NewDebt    string `json:"new_debt"`
	PaymentId  string `json:"payment_id"`
}
