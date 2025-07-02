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
	"oprosdom.ru/monolith/internal/dz/internal/handlers"
)

func main() {

	// предусмотреть контекст!
	// 1) если клиент стопнул в браузере выполнение, то нужно отменять операции -> это предусмотрено http сервером, но нужно обрабатывать  это событие в хендлерах
	// 2) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/api/member", handlers.AddMember).Methods("POST")
	r.HandleFunc("/api/kvartira", handlers.AddKvartira).Methods("POST")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
		// таймауты read/write тут не указываем потому что у нас вероятно будут вебсокеты на этом же порту, поэтому будем другими механизмами таймауты отслеживать
		// хотя нужно изучить вопрос, ws это же не http, а поверх http
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
