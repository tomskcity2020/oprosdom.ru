package repo_internal

import (
	"log"

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

func (repo *RepositoryStruct) Save(m models.ModelInterface) {

	//fmt.Println("in save now")
	switch m.Type() {
	case "member":
		if member, ok := m.(*models.Member); ok {
			repo.members = append(repo.members, member)
		} else {
			log.Println("Проблема с приведением к типу Member")
		}

	case "kvartira":
		if kvartira, ok := m.(*models.Kvartira); ok {
			repo.kvartiras = append(repo.kvartiras, kvartira)
		} else {
			log.Println("Проблема с приведением к типу Kvartira")
		}
	}
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
