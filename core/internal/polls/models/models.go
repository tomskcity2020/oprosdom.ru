package polls_models

type Poll struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type PollStats struct {
	PollID     int `json:"poll_id"`
	ZaCount    int `json:"za_count"`
	ProtivCount int `json:"protiv_count"`
}