package handlers

import (
	"encoding/json"
	"net/http"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	//"oprosdom.ru/monolith/internal/dz/internal/service"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME"))
}

func AddMember(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// парсим
	var member models.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// создаем uuid (отдельно делаем чтоб BasicValidate использовать и для POST и для PUT)
	if err := member.CreateUuid(); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
	}

	// проводим первичную проверку (структура заполнена полностью сейчас)
	if err := member.BasicValidate(); err != nil {
		replyError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// сохраняем в файл
	if err := repo.GetRepoSingleton().SaveToFile(&member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// сохраняем в слайс
	if err := repo.GetRepoSingleton().Save(&member); err != nil {
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

func AddKvartira(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// парсим
	var kvartira models.Kvartira
	if err := json.NewDecoder(r.Body).Decode(&kvartira); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// создаем uuid (отдельно делаем чтоб BasicValidate использовать и для POST и для PUT)
	if err := kvartira.CreateUuid(); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
	}

	// проводим первичную проверку (структура заполнена полностью сейчас)
	if err := kvartira.BasicValidate(); err != nil {
		replyError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// сохраняем в файл
	if err := repo.GetRepoSingleton().SaveToFile(&kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// сохраняем в слайс
	if err := repo.GetRepoSingleton().Save(&kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&kvartira)

}
