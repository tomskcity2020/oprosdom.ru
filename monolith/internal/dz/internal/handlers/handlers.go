package handlers

import (
	"log"
	"net/http"

	"oprosdom.ru/monolith/internal/dz/internal/models"
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
