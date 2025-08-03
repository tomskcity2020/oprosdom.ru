package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"oprosdom.ru/microservice_notify/internal/transport/kafka"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// userRepoConn := "postgres://test:test@127.0.0.1:5432/users?" +
	// 	"sslmode=disable&" +
	// 	"pool_min_conns=5&" +
	// 	"pool_max_conns=25&" +
	// 	"pool_max_conn_lifetime=30m&" +
	// 	"pool_max_conn_lifetime_jitter=5m&" +
	// 	"pool_max_conn_idle_time=15m&" +
	// 	"pool_health_check_period=1m"

	// usersRepo, err := repo.NewRepoFactory(ctx, userRepoConn)
	// if err != nil {
	// 	log.Fatalf("repo initialization failed with error: %v", err)
	// }
	// defer usersRepo.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	// codeTransport, err := transport.NewTransportFactory(ctx, "localhost:9092", "code")
	// if err != nil {
	// 	log.Fatalf("codeTransport initialization failed with error: %v", err)
	// }
	// defer codeTransport.Close()

	// usersService := service.NewServiceFactory(usersRepo, codeTransport)
	// usersHandler := handlers.NewHandler(usersService)

	// предусмотреть контекст!
	// 1) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались. Чтоб не получилось так, что из кафки возьмем сообщение и завершимся. Нужно обязательно чтоб отправка происходила.

	// в случае с http сервером мы загоняли сервис в хендлер, здесь же мы загоняем сервис в kafkaConsumer
	//kafkaConsumer := kafka.NewConsumer([]string{"localhost:9092"}, cfg.KafkaTopic, cfg.KafkaGroupID, svc)
	kafkaConsumer := kafka.NewConsumer([]string{"localhost:9092"}, "code", "for_notify")

	errCh := make(chan error, 1)

	go func() {
		if err := kafkaConsumer.Start(ctx); err != nil {
			errCh <- err
		}
	}()

	// select блокирует программу до наступления одного из case. бесконечный for здесь смысла не имеет, так как select запускается единожды и результат любого case приводит к завершению программы
	select {
	case err := <-errCh: // присваиваем err значение из канала, область видимости только этот case
		log.Fatalf("Возникла ошибка: %v", err)
	case <-ctx.Done():
		// Даем разумное время на завершение обработки отправки sms
		// ctx контекст не передаем иначе shutdownCtx немедленно отменится при отмене ctx. Нужно делать новый контекст чтоб дать время
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Здесь можно добавить остановку других компонентов, если они есть

		// Ожидаем завершение shutdownCtx
		<-shutdownCtx.Done()
		log.Println("Graceful Shutdown успешен")
	}

}
