package repo_internal

import (
	"oprosdom.ru/monolith/internal/dz6/internal/models"
)

type RepositoryStruct struct {
	members   []models.ModelInterface
	kvartiras []models.ModelInterface
}

func NewCallInternalRepo() *RepositoryStruct {
	return &RepositoryStruct{
		members:   make([]models.ModelInterface, 0),
		kvartiras: make([]models.ModelInterface, 0),
	}
}

func (repo *RepositoryStruct) Save(m models.ModelInterface) {

	//fmt.Println("in save now")
	switch m.Type() {
	case "member":
		repo.members = append(repo.members, m)
	case "kvartira":
		repo.kvartiras = append(repo.kvartiras, m)
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
