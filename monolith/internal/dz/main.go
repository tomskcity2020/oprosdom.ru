package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// логика graceful shutdown у нас такая: дожидаемся записи в слайсы и прекращаем дальнейшее выполнение цикла, выходим. Потому что посреди записи в слайс нелогично делать graceful shutdown, потому что он таковым являться не будет ввиду того, что часть данных запишется в слайс, а часть возможно нет
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {

		select {
		case <-ctx.Done():
			log.Println("GRACEFUL SHUTDOWN main")
			return
		default:
			serviceEntity := service.NewServiceFactory()

			startTime := time.Now()

			serviceEntity.RunParallel(ctx, modelsData)
			//serviceEntity.RunSeq(modelsData)

			time.Sleep(1 * time.Second) // имитация длит выполнения

			elapsedTime := time.Since(startTime)
			log.Printf("Функция Run завершилась за: %v\n", elapsedTime)

		}

	}

}
