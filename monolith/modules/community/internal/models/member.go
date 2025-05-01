package models

type Member_personal struct {
	id              int
	Firstname       string
	Lastname        string
	mobile          string
	mobile_verified bool
	nickname        string
	status          string // active limited spamban fake - если чела в одном из сообществ забанили за спам, то нужно spamban статус применять к базовой структуре, чтоб во всех сообществах накладывать ограничения
}

// участник может состоять в нескольких домах, но так как верификация везде своя (старшим или сообществом), то разделяем персональные данные и данные по Жителю

type Member struct {
	Member_personal
	community_id             int
	kv                       string
	verified_by_card         bool
	verified_by_housemanager bool
	verified_by_community    bool
	temporary_ban            bool // чел может быть зареган в нескольких Сообществах, в одном что-то натворить и схлопотать бан, это значит что бан распространяется только на this сообщество
}
