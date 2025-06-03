package repo_internal

import (
	"log"
	"sync"
	"time"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

type RepositoryStruct struct {
	members   []*models.Member
	kvartiras []*models.Kvartira
}

func NewCallInternalRepo() *RepositoryStruct {
	return &RepositoryStruct{
		members:   make([]*models.Member, 0),
		kvartiras: make([]*models.Kvartira, 0),
	}
}

// func (repo *RepositoryStruct) Save(ch <-chan models.ModelInterface) {
func (repo *RepositoryStruct) Save(m models.ModelInterface) {

	// сюда передаем только интерфейс
	// убираем в сервисный слой вызов горутин
	// здесь оставляем только приведение к типам и запись в слайсы
	// но мьютексы оставляем здесь, так как писать будем в слайсы из нескольких горутин (которые будем вызывать в сервисном слое)

	muMembers := sync.Mutex{}
	muKvartiras := sync.Mutex{}

	switch m.Type() {
	case "member":
		if member, ok := m.(*models.Member); ok {
			muMembers.Lock()
			repo.members = append(repo.members, member)
			muMembers.Unlock()
		} else {
			log.Println("Проблема с приведением к типу Member")
		}

	case "kvartira":
		if kvartira, ok := m.(*models.Kvartira); ok {
			muKvartiras.Lock()
			repo.kvartiras = append(repo.kvartiras, kvartira)
			muKvartiras.Unlock()
		} else {
			log.Println("Проблема с приведением к типу Kvartira")
		}
	default:
		log.Println("Неведомый тип")
	}

}

func (repo *RepositoryStruct) Check() {

	go func() {
		// узнавать количество структур в каждом слайсе на старте (однозначно там по 0)
		prevMembersCount := len(repo.members)
		prevKvartirasCount := len(repo.kvartiras)

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop() // обязательно освобождаем ресурсы, сработает когда основная горутина завершится

		for range ticker.C {
			// узнавать количество структур в каждом слайсе каждые 200 мс
			nowMembersCount := len(repo.members)
			nowKvartirasCount := len(repo.kvartiras)

			if nowMembersCount > prevMembersCount {
				newMembers := repo.members[prevMembersCount:] // все новые это те, индекс которых идет следом за prevMembersCount
				for _, member := range newMembers {
					log.Printf("New value in Members: %+v ", member)
				}
				prevMembersCount = nowMembersCount
			}

			if nowKvartirasCount > prevKvartirasCount {
				newKvartiras := repo.kvartiras[prevKvartirasCount:]
				for _, kvartira := range newKvartiras {
					log.Printf("New value in Kvartiras: %+v ", kvartira)
				}
				prevKvartirasCount = nowKvartirasCount
			}

		}
	}()

}

func (repo *RepositoryStruct) Show(t string) int {

	//fmt.Println("in show now")
	switch t {
	case "member":
		return len(repo.members)
	case "kvartira":
		return len(repo.kvartiras)
	default:
		return 0
	}
}
