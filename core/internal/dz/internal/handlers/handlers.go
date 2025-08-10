package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"oprosdom.ru/core/internal/dz/internal/models"
	"oprosdom.ru/core/internal/dz/internal/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME"))
}

func (h *Handler) MemberAdd(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// парсим
	var member models.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// 

	// service (там репо и бизнес)
	if err := h.service.MemberAdd(r.Context(), &member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&member)

}

func replyError(w http.ResponseWriter, error string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": error,
	})
}

func (h *Handler) KvartiraAdd(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// парсим
	var kvartira models.Kvartira
	if err := json.NewDecoder(r.Body).Decode(&kvartira); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// service (там репо и бизнес)
	if err := h.service.KvartiraAdd(r.Context(), &kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&kvartira)

}

func (h *Handler) MemberUpdate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var member models.Member

	// извлекаем id, если нет id, то никакой ошибки или паники не будет, получим пустую строку ""
	varsMap := mux.Vars(r)
	idRaw := varsMap["id"]

	// добавляем пришедший id в структуру, пока не проверяем - проверка будет далее всех полей структуры сразу
	member.Id = idRaw

	// парсим данные в структуру
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// service (там репо и бизнес)
	if err := h.service.MemberUpdate(r.Context(), &member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&member)

}

func (h *Handler) KvartiraUpdate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var kvartira models.Kvartira

	// извлекаем id, если нет id, то никакой ошибки или паники не будет, получим пустую строку ""
	varsMap := mux.Vars(r)
	idRaw := varsMap["id"]

	// добавляем пришедший id в структуру, пока не проверяем - проверка будет далее всех полей структуры сразу
	kvartira.Id = idRaw

	// парсим данные в структуру
	if err := json.NewDecoder(r.Body).Decode(&kvartira); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// service (там репо и бизнес)
	if err := h.service.KvartiraUpdate(r.Context(), &kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&kvartira)

}

func (h *Handler) MembersGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// service (там репо и бизнес)
	data, err := h.service.MembersGet(r.Context())
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)

}

func (h *Handler) KvartirasGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// service (там репо и бизнес)
	data, err := h.service.KvartirasGet(r.Context())
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)

}

func (h *Handler) MemberGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	varsMap := mux.Vars(r)
	id := varsMap["id"]

	// service (там репо и бизнес)
	data, err := h.service.MemberGet(r.Context(), id)
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)

}

func (h *Handler) KvartiraGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	varsMap := mux.Vars(r)
	id := varsMap["id"]

	// service (там репо и бизнес)
	data, err := h.service.KvartiraGet(r.Context(), id)
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)

}

func (h *Handler) RemoveById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	varsMap := mux.Vars(r)
	mk := varsMap["mk"]
	id := varsMap["id"]

	// service (там репо и бизнес)
	if err := h.service.RemoveById(r.Context(), id, mk); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PayDebt(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var request models.PayDebtRequest

	varsMap := mux.Vars(r)
	request.MemberId = varsMap["id"]


	// парсим данные в структуру то что пришло по POST
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// service (там репо и бизнес)
	data, err := h.service.PayDebt(r.Context(), &request)

	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&data)


}