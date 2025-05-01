package models

type Community struct {
	id            int
	houseguid     string // (фиас uniq)
	country       string
	region        string
	naspunkt      string
	address       string
	latitude      float64
	longitude     float64
	started_epoch int64
	users_count   int
	status        string // waiting5 / pending30 / success
}

// методы:
// ShowCommunity
// AddCommunity
