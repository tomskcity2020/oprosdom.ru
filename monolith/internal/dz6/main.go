package main

import (
	"oprosdom.ru/monolith/internal/dz6/internal/service"
)

func main() {

	service := service.NewServiceFactory()
	service.Run()

}
