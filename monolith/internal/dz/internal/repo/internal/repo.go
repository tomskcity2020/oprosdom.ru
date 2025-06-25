package repo_internal

import (
	"encoding/json"
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
}

func NewCallInternalRepo() *RepositoryStruct {
	return &RepositoryStruct{
		members:   make([]*models.Member, 0),
		kvartiras: make([]*models.Kvartira, 0),
		// muMembers и muKvartiras автоматически инициализируются
	}

}

func (repo *RepositoryStruct) Save(m models.ModelInterface) {

	//switch m.Type() {
	switch data := m.(type) {
	case *models.Member:
		//time.Sleep(300 * time.Millisecond) // слип для эмуляции времени работы например записи в базу данных или отправки данных через grpc
		repo.muMembers.Lock()
		defer repo.muMembers.Unlock()
		//log.Printf("data: %+v", data)
		repo.members = append(repo.members, data)
		//log.Printf("repo.members: %+v", repo.members)

	case *models.Kvartira:
		//time.Sleep(300 * time.Millisecond) // слип для эмуляции времени работы например записи в базу данных или отправки данных через grpc
		repo.muKvartiras.Lock()
		defer repo.muKvartiras.Unlock()
		repo.kvartiras = append(repo.kvartiras, data)
		//log.Println("repo add kvartiras done")
	default:
		log.Println("Неведомый тип")
	}

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

func (repo *RepositoryStruct) SaveToFile(fileName string) {
	// мьютекс использовать здесь излишне, так как запись в файл будет осуществляться не параллельно
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("ошибка создания %v: %v", fileName, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	switch fileName {
	case "members.json":
		if err := encoder.Encode(repo.members); err != nil {
			log.Printf("ошибка записи в %v: %v", fileName, err)
		}
	case "kvartiras.json":
		if err := encoder.Encode(repo.kvartiras); err != nil {
			log.Printf("ошибка записи в %v: %v", fileName, err)
		}
	}

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
