package users_dto

// собрали dto в отдельный пакет для того, чтобы к нему могли обращаться и из handler'ов и из сервисного слоя. Потому как если разместить dto, например, в обработчике, то может возникнуть перекрестный вызов: из обработчика ссылаемся на сервисный слой, а сервисный слой ссылается на обработчик

type RequestSignUpDTO struct {
	Phone     string `json:"phone"`
	City      string `json:"city"`
	HouseId   string `json:"house_id"`
	AptNumber string `json:"apt_number"`
	Number    int    `json:"number"`
}

type ResponseSignUpDTO struct {
	PublicId int `json:"public_id"`
}
