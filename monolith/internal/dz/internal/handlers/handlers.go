package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	"oprosdom.ru/monolith/internal/dz/internal/service"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME"))

	member := models.Member{}
	log.Printf("%v", member)

	kvartira := models.Kvartira{}
	log.Printf("%v", kvartira)

	service := service.NewServiceFactory()
	service.CountData()
}

func AddMember(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// парсим
	var member models.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// проводим первичную проверку
	if err := member.BasicValidate(); err != nil {
		replyError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// сохраняем в файл и слайс (модернизируем репо: добавить возврат ошибки)
	if err := repo.GetRepoSingleton().Save(&member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)

}

func replyError(w http.ResponseWriter, error string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": error,
	})
}
