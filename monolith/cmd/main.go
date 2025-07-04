package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	core "oprosdom.ru/monolith/internal"
)

func main() {

	// 	Вкратце схема такая:
	// 1) вызываем rpc обработчик общий
	// 2) он в зависимости от json-rpc метода выбирает спец обработчик и назначает ему соответствующий части приложения репозиторий
	// 3) далее спец обработчик парсит и создает дто и проверяем
	// 4) отдаем дто в сервисный слой
	// 5) сервисный слой вызывает репозиторий и бизнес-логику

	rpc_handler := core.NewJsonRpcHandler()

	r := mux.NewRouter()
	r.HandleFunc("/rpc", rpc_handler.RequestResponse).Methods("POST")

	log.Println("JSON-RPC сервер запущен на порте 8080.")
	log.Fatal(http.ListenAndServe(":8080", r))

}
