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
	"oprosdom.ru/core/internal/dz/internal/handlers"
	"oprosdom.ru/core/internal/dz/internal/repo"
	"oprosdom.ru/core/internal/dz/internal/service"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo, err := repo.NewRepoFactory(ctx, "postgres://test:test@127.0.0.1:5432/test?sslmode=disable")
	if err != nil {
		log.Fatalf("repo initialization failed with error: %v", err)
	}
	defer repo.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	service := service.NewServiceFactory(repo)
	h := handlers.NewHandler(service)

	// предусмотреть контекст!
	// 1) если клиент стопнул в браузере выполнение, то нужно отменять операции -> это предусмотрено http сервером, но нужно обрабатывать  это событие в хендлерах
	// 2) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались

	// curl -X POST "http://127.0.0.1:8080/api/member" -H "Content-Type: application/json" -d '{"name":"Иван Иванов","phone":"+79991234567","community":5}'
	// curl -X POST "http://127.0.0.1:8080/api/kvartira" -H "Content-Type: application/json" -d '{"number":"115","komnat":2}'
	// curl -X PUT "http://127.0.0.1:8080/api/member/09770685-ae8a-4e68-9751-50bba1d846f1" -H "Content-Type: application/json" -d '{"name":"Иван Иванов","phone":"+79991234567","community":5}'
	// curl -X PUT "http://127.0.0.1:8080/api/kvartira/1f970b7b-679c-4c7d-a252-3ef370d439f4" -H "Content-Type: application/json" -d '{"number":"12","komnat":3}'
	// curl -X GET "http://127.0.0.1:8080/api/members" -H "Content-Type: application/json"
	// curl -X GET "http://127.0.0.1:8080/api/kvartiras" -H "Content-Type: application/json"
	// curl -X GET "http://127.0.0.1:8080/api/member/09770685-ae8a-4e68-9751-50bba1d846f1"
	// curl -X GET "http://127.0.0.1:8080/api/kvartira/1f970b7b-679c-4c7d-a252-3ef370d439f4"
	// curl -X DELETE "http://127.0.0.1:8080/api/member/09770685-ae8a-4e68-9751-50bba1d846f1"
	// curl -X DELETE "http://127.0.0.1:8080/api/kvartira/1f970b7b-679c-4c7d-a252-3ef370d439f4"
	// curl -X POST "http://127.0.0.1:8080/api/member/b406690d-018a-4fb5-b1b7-70946b92432e/paydebt" -H "Content-Type: application/json" -d '{"kvartira_id":"f983de63-e4f3-483a-b9a2-6e1f3ea960b7","amount":"17500.55"}'

	r := mux.NewRouter()
	r.HandleFunc("/", h.HomeHandler)
	r.HandleFunc("/api/member", h.MemberAdd).Methods("POST")
	r.HandleFunc("/api/kvartira", h.KvartiraAdd).Methods("POST")
	r.HandleFunc("/api/member/{id}", h.MemberUpdate).Methods("PUT")
	r.HandleFunc("/api/kvartira/{id}", h.KvartiraUpdate).Methods("PUT")
	r.HandleFunc("/api/members", h.MembersGet).Methods("GET")
	r.HandleFunc("/api/kvartiras", h.KvartirasGet).Methods("GET")
	r.HandleFunc("/api/member/{id}", h.MemberGet).Methods("GET")
	r.HandleFunc("/api/kvartira/{id}", h.KvartiraGet).Methods("GET")
	r.HandleFunc("/api/{mk}/{id}", h.RemoveById).Methods("DELETE")
	r.HandleFunc("/api/member/{id}/paydebt", h.PayDebt).Methods("POST")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
		// таймауты read/write тут не указываем потому что у нас вероятно будут вебсокеты на этом же порту, поэтому будем другими механизмами таймауты отслеживать
		// хотя нужно изучить вопрос, ws это же не http, а поверх http
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("Стартуем noTLS сервак на %v порту", srv.Addr)
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
		shutdownCtx, shutdownCtxStop := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCtxStop()

		// Shutdown() перестанет принимать новые подключения, но завершит старые
		// на заметку Shutdown() не ждёт завершения WS-соединений
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful Shutdown http сервера не сработало, ошибка: %v", err)
		}
		log.Println("Graceful Shutdown http сервера успешен")
	}

}
