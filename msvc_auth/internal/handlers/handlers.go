package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/msvc_auth/internal/service"
	"oprosdom.ru/shared"
)

type Handler struct {
	service service.UsersService
}

func NewHandler(service service.UsersService) *Handler {
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

	validatedCodeCheckReq, err := unsafeCodeCheckReq.Validate()
	if err != nil {
		replyError(w, err, "incorrect_phonesend", http.StatusBadRequest)
		return
	}

	if err := h.service.CodeCheck(r.Context(), validatedCodeCheckReq); err != nil {
		replyError(w, err, "codecheck_wrong_request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("success")

}
