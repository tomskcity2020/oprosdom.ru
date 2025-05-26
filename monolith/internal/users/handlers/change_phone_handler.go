package handlers

type RequestChangePhone struct {
	Phone string `json:"phone"`
}

type ResponseChangePhone struct {
	FormattedNewPhone string `json:"new_phone"`
}

type ChangePhone struct{}

func (obj *ChangePhone) Handle(params map[string]interface{}) (interface{}, error) {
	// парсим
	// валидируем
	// вызываем сервисный слой соответствующий (далее в нем уже вызывается бизнес логика)
	// получаем ответ от сервисного слоя и на его основе создаем ответное дто
	// отдаем dto (далее в core мы его маршалим в json)
	result := "new phone"
	return &result, nil
}
