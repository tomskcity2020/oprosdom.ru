package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/msvc_auth/internal/service"
	"oprosdom.ru/shared"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func replyError(w http.ResponseWriter, err error, publicErr string, statusCode int) {
	// оригинальный err.Error() логируем в redis с ttl
	// TODO
	log.Println(err.Error())
	w.WriteHeader(statusCode)
	// отправлять error поле в json'e нет смысла, так как ошибка будет по status code определяться
	json.NewEncoder(w).Encode(publicErr)
}

func (h *Handler) PhoneSend(w http.ResponseWriter, r *http.Request) {

	// 1) парсим json
	// 2) проводим санитизацию и первичную валидацию
	// 3) формируем структуру и отдаем в сервис
	// 4) ошибки возвращаем: err логируем, publicErr отдаем фронтенду. publicErr не должен содержать чувствительных данных.

	w.Header().Set("Content-Type", "application/json")

	// парсим и наполняем unsafe структуру данными
	var unsafePhoneSendReq models.UnsafePhoneSendReq
	if err := json.NewDecoder(r.Body).Decode(&unsafePhoneSendReq); err != nil {
		replyError(w, err, "incorrect_request", http.StatusBadRequest)
		return
	}

	unsafePhoneSendReq.UserAgent = r.UserAgent()
	unsafePhoneSendReq.Ip = shared.IpHttpGet(r)

	// проводим первичную валидацию, на выходе заполняем valid структуру и отдаем в сервисный слой (для дальнейшей бизнес-валидации, бизнес-логики и взаимодействия с репо)
	validatedPhoneSendReq, err := unsafePhoneSendReq.Validate()
	if err != nil {
		replyError(w, err, "incorrect_phonesend", http.StatusBadRequest)
		return
	}

	// service (там репо и бизнес)
	if err := h.service.PhoneSend(r.Context(), validatedPhoneSendReq); err != nil {
		replyError(w, err, "phonesend_wrong_request", http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("success")

}

func (h *Handler) CodeCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var unsafeCodeCheckReq models.UnsafeCodeCheckReq
	if err := json.NewDecoder(r.Body).Decode(&unsafeCodeCheckReq); err != nil {
		replyError(w, err, "incorrect_request", http.StatusBadRequest)
		return
	}

	unsafeCodeCheckReq.UserAgent = r.UserAgent()
	unsafeCodeCheckReq.Ip = shared.IpHttpGet(r)

	validatedCodeCheckReq, err := unsafeCodeCheckReq.Validate()
	if err != nil {
		replyError(w, err, "codecheck_validation_failed", http.StatusBadRequest)
		return
	}

	if err := h.service.CodeCheck(r.Context(), validatedCodeCheckReq); err != nil {
		replyError(w, err, "codecheck_wrong_request", http.StatusInternalServerError)
		return
	}

	// 10 лет в часах (для jwt и куки)
	tenYears := time.Hour * 24 * 365 * 10

	jwtTokenStr, err := h.service.CreateJwt(r.Context(), tenYears, validatedCodeCheckReq)
	if err != nil {
		replyError(w, err, "create_token_failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    jwtTokenStr,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().AddDate(10, 0, 0), // 10 лет
		MaxAge:   86400 * 365 * 10,             // 10 лет
	})

	// удаляем код, иначе могут 1 успешным кодом создать очень много токенов
	if err := h.service.PurgeCode(r.Context(), validatedCodeCheckReq); err != nil {
		// просто логируем, не прерываем выполнение
		log.Printf("purgecode failed: %v", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "Token issued in cookie"}`))

}
