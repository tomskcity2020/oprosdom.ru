package repo_internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

type RepositoryStruct struct {
	// добавляем мьютекс в структуру для того, чтобы в каждом методе не создавать их - так они потеряют свой смысл (потому что параллельно запускаем репо). А так мы изначально создаем репо через конструктор и мьютексы используются во всех горутинах
	muMembers   sync.RWMutex
	members     []*models.Member
	muKvartiras sync.RWMutex
	kvartiras   []*models.Kvartira
	muRwFile    sync.RWMutex
}

func NewCallInternalRepo() *RepositoryStruct {
	return &RepositoryStruct{
		members:   make([]*models.Member, 0),
		kvartiras: make([]*models.Kvartira, 0),
		// muMembers и muKvartiras автоматически инициализируются
	}

}

func (repo *RepositoryStruct) Save(m models.ModelInterface) error {

	//switch m.Type() {
	switch data := m.(type) {
	case *models.Member:
		//time.Sleep(300 * time.Millisecond) // слип для эмуляции времени работы например записи в базу данных или отправки данных через grpc
		repo.muMembers.Lock()
		defer repo.muMembers.Unlock()
		//log.Printf("data: %+v", data)
		repo.members = append(repo.members, data)
		jsonData, _ := json.Marshal(repo.members)
		log.Printf("[]members: %s", jsonData)

	case *models.Kvartira:
		//time.Sleep(300 * time.Millisecond) // слип для эмуляции времени работы например записи в базу данных или отправки данных через grpc
		repo.muKvartiras.Lock()
		defer repo.muKvartiras.Unlock()
		repo.kvartiras = append(repo.kvartiras, data)
		//log.Println("repo add kvartiras done")
		jsonData, _ := json.Marshal(repo.kvartiras)
		log.Printf("[]kvartiras: %s", jsonData)
	default:
		return errors.New("неведомый тип save")
	}

	return nil

}

func (repo *RepositoryStruct) LoadFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// если файла нет, это значит первый запуск, просто выходим
			return
		}
		log.Printf("Файл есть, но чтение не удалось %v: %v", fileName, err)
		return
	}
	defer file.Close()

	switch fileName {
	case "members.json":
		if err := json.NewDecoder(file).Decode(&repo.members); err != nil {
			log.Printf("Некорректный формат файла %v: %v", fileName, err)
		}
	case "kvartiras.json":
		if err := json.NewDecoder(file).Decode(&repo.kvartiras); err != nil {
			log.Printf("Некорректный формат файла %v: %v", fileName, err)
		}
	}

}

func (repo *RepositoryStruct) MembersInSliceNow() int {
	return len(repo.members)
}

func (repo *RepositoryStruct) KvartirasInSliceNow() int {
	return len(repo.kvartiras)
}

func (repo *RepositoryStruct) SaveToFile(m models.ModelInterface) error {

	repo.muRwFile.Lock()
	defer repo.muRwFile.Unlock()

	filename := ""
	var recData any

	switch data := m.(type) {
	case *models.Member:
		filename = "members.json"
		recData = data
	case *models.Kvartira:
		filename = "kvartiras.json"
		recData = data
	default:
		return errors.New("неведомый тип savetofile")
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла %v для записи: %v", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(recData); err != nil {
		return fmt.Errorf("ошибка записи в файл %v: %v", filename, err)
	}

	return nil

}

func (repo *RepositoryStruct) GetSliceMembers() []*models.Member {
	// для чтения используем спец мьютекс RWMutex: читать могут параллельно из множества горутин, если запись - никто не может читать
	// TODO нужно посмотреть вероятно обычный мьютекс быстрее, иначе бы везде использовали RW. Пока у нас только тикер использует этот метод и возможно RW тут лишний
	repo.muMembers.RLock()
	defer repo.muMembers.RUnlock()
	return repo.members
}

func (repo *RepositoryStruct) GetSliceKvartiras() []*models.Kvartira {
	repo.muKvartiras.RLock()
	defer repo.muKvartiras.RUnlock()
	return repo.kvartiras
}
