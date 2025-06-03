package main

import (
	"oprosdom.ru/monolith/internal/dz/internal/service"
)

func main() {

	serviceEntity := service.NewServiceFactory()

	serviceEntity.Run()

}
