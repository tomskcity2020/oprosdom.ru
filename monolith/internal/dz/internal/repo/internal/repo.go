package repo_internal

import (
	"log"
	"sync"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

type RepositoryStruct struct {
	members   []*models.Member
	kvartiras []*models.Kvartira
	// добавляем мьютекс в структуру для того, чтобы в каждом методе не создавать их - так они потеряют свой смысл (потому что параллельно запускаем репо). А так мы изначально создаем репо через конструктор и мьютексы используются во всех горутинах
	muMembers   sync.RWMutex
	muKvartiras sync.RWMutex
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
		repo.members = append(repo.members, data)
		//log.Println("repo add members done")

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
