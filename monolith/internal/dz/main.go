package main

import (
	"oprosdom.ru/monolith/internal/dz/internal/service"
)

func main() {

	service := service.NewServiceFactory()
	service.Run()

}
