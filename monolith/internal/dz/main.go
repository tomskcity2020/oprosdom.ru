package main

import (
	"fmt"
	"time"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/service"
)

func main() {

	modelsData := []models.ModelInterface{
		models.NewUserFactory("Namme", "+79991231234", 54),
		models.NewUserFactory("Bobby", "+71000033214", 79),
		models.NewKvartiraFactory("135b", 3),
		models.NewUserFactory("Marry", "+72223342343", 91),
		models.NewKvartiraFactory("179", 1),
		models.NewKvartiraFactory("11a", 2),
		models.NewUserFactory("Alex", "+71000032344", 51),
	}

	//for i := 0; i < 10; i++ {

	serviceEntity := service.NewServiceFactory()

	startTime := time.Now()

	serviceEntity.RunParallel(modelsData)
	//serviceEntity.RunSeq(modelsData)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Функция Run завершилась за: %v\n", elapsedTime)

	time.Sleep(time.Second) // для демонстрации работы функции-чекера слайсов при RunParallel иначе все исполнится быстрее, чем вторая итерация чекера наступит

	//}

}
