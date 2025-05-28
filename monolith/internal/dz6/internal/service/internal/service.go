package service_internal

import (
	"fmt"
	"time"

	"oprosdom.ru/monolith/internal/dz6/internal/models"
	"oprosdom.ru/monolith/internal/dz6/internal/repo"
)

type ServiceStruct struct {
	repo repo.RepositoryInterface
}

func NewCallInternalService(repo repo.RepositoryInterface) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
	}
}

func (obj *ServiceStruct) Run() {

	fmt.Println("Сервис запущен. Пожалуйста, ожидайте..")

	models := []models.ModelInterface{
		models.NewUserFactory("Namme", "+79991231234", 54),
		models.NewUserFactory("Bobby", "+71000033214", 79),
		models.NewKvartiraFactory("135b", 3),
		models.NewUserFactory("Marry", "+72223342343", 91),
		models.NewKvartiraFactory("179", 1),
		models.NewKvartiraFactory("11a", 2),
		models.NewUserFactory("Alex", "+71000032344", 51),
	}

	totalElements := len(models)

	for i, model := range models {
		obj.repo.Save(model)
		time.Sleep(2 * time.Second)
		fmt.Printf("Обработано: %v из %v", i+1, totalElements)
		fmt.Println()
	}

	showMembers := obj.repo.Show("member")
	fmt.Printf("Всего участников: %v", showMembers)
	fmt.Println()

	showKvartiras := obj.repo.Show("kvartira")
	fmt.Printf("Всего квартир: %v", showKvartiras)
	fmt.Println()

}
