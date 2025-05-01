package main

import (
	"fmt"

	"oprosdom.ru/monolith/cmd/internal/model"
)

func main() {

	member := model.Member_personal{}
	member.Set_Firstname("Сергей")
	fmt.Println("Firstname:", member.Firstname())

}
