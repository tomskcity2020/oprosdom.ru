package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"oprosdom.ru/microservice_auth/internal/handlers"
	"oprosdom.ru/microservice_auth/internal/repo"
	"oprosdom.ru/microservice_auth/internal/service"
	"oprosdom.ru/microservice_auth/internal/transport"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	userRepoConn := "postgres://test:test@127.0.0.1:5433/users?" +
		"sslmode=disable&" +
		"pool_min_conns=5&" +
		"pool_max_conns=25&" +
		"pool_max_conn_lifetime=30m&" +
		"pool_max_conn_lifetime_jitter=5m&" +
		"pool_max_conn_idle_time=15m&" +
		"pool_health_check_period=1m"

	usersRepo, err := repo.NewRepoFactory(ctx, userRepoConn)
	if err != nil {
		log.Fatalf("repo initialization failed with error: %v", err)
	}
	defer usersRepo.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	codeTransport, err := transport.NewTransportFactory(ctx, "localhost:9092", "code")
	if err != nil {
		log.Fatalf("codeTransport initialization failed with error: %v", err)
	}
	defer codeTransport.Close()

	usersService := service.NewServiceFactory(usersRepo, codeTransport)
	usersHandler := handlers.NewHandler(usersService)

	// предусмотреть контекст!
	// 1) если клиент стопнул в браузере выполнение, то нужно отменять операции -> это предусмотрено http сервером, но нужно обрабатывать  это событие в хендлерах
	// 2) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались

	// curl -X POST "http://127.0.0.1/auth/requestcode" -H "Content-Type: application/json" -d '{"phone":"+79991234567"}'

	r := mux.NewRouter()
	r.HandleFunc("/auth/requestcode", usersHandler.PhoneSend).Methods("POST")
	// verification/phonecode

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
		// таймауты read/write тут не указываем потому что у нас вероятно будут вебсокеты на этом же порту, поэтому будем другими механизмами таймауты отслеживать
		// хотя нужно изучить вопрос, ws это же не http, а поверх http
	}

	errCh := make(chan error, 1)

	// запускаем в отдельной горутине потому что ListenAndServe это блокирующий вызов, и без горутины мы до select никогдай не дойдем. А значит graceful shutdown не получится.
	go func() {
		log.Printf("Стартуем noTLS микросервис AUTH на %v порту", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	// select блокирует программу до наступления одного из case. бесконечный for здесь смысла не имеет, так как select запускается единожды и результат любого case приводит к завершению программы
	select {
	case err := <-errCh: // присваиваем err значение из канала, область видимости только этот case
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Println("Сервак остановлен. Это не неожиданная ошибка, а ожидаемое поведение завершения")
		default:
			log.Fatalf("Возникла ошибка: %v", err)
		}
	case <-ctx.Done():
		// мы не должны передавать ctx иначе shutdownCtx немедленно отменится при отмене ctx. Вместо этого делаем новый контекст, который даст 5 сек на довыполнение
		shutdownCtx, shutdownCtxStop := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCtxStop()

		// Shutdown() перестанет принимать новые подключения, но завершит старые
		// на заметку Shutdown() не ждёт завершения WS-соединений
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful Shutdown http сервера не сработало, ошибка: %v", err)
		}
		log.Println("Graceful Shutdown http сервера успешен")
	}

}
