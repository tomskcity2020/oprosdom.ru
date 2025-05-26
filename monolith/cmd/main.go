package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	core "oprosdom.ru/monolith/internal"
	"oprosdom.ru/monolith/internal/dz6"
)

func dz6run() {
	models := []string{"member", "kvartira"}

	x := 0
	for {

		for i := 0; i <= 1; i++ {

			modelStruct, err := dz6.CreateModel(models[i])
			if err != nil {
				fmt.Println(err)
			}

			result, err := modelStruct.ToSlice()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%+v", result)
			fmt.Println()

			time.Sleep(1 * time.Second)

		}

		x++
		if x > 3 {
			break
		}
	}

}

func main() {

	// запускаем домашку в отдельной горутине
	// делаем слайс типов, которые будем рандомом передавать с помощью первой функции во вторую
	// запускаем это в бесконечном цикле со слипом в 1 секунду

	go dz6run()

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
