package repo_internal

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strconv"
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

func NewCallInternalRepo() (*RepositoryStruct, error) {
	repo := &RepositoryStruct{
		members:   make([]*models.Member, 0),
		kvartiras: make([]*models.Kvartira, 0),
		// muMembers и muKvartiras автоматически инициализируются
	}

	// заполняем слайсы данными из файлов
	if err := repo.loadMembersFromFile(); err != nil {
		return nil, err
	}

	if err := repo.loadKvartirasFromFile(); err != nil {
		return nil, err
	}

	return repo, nil

}

func (repo *RepositoryStruct) loadMembersFromFile() error {
	// один мьютекс охватывает и файл и слайс
	repo.muMembers.Lock()
	defer repo.muMembers.Unlock()

	file, err := os.Open("members.csv")
	if err != nil {
		// если файла не существует значит первый запуск вероятнее всего, а значит дальнейшие действия по парсингу не имеют смысла
		if errors.Is(err, fs.ErrNotExist) { // os.IsNotExist старый метод, его не юзаем
			log.Println("файл members.csv не существует, первый запуск?")
			return nil
		}
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	count := 0

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		community, _ := strconv.Atoi(record[3])

		member := &models.Member{
			Id:        record[0],
			Name:      record[1],
			Phone:     record[2],
			Community: community,
		}

		repo.members = append(repo.members, member)
		count++
	}

	log.Printf("%v жителей загружено из файла", count)

	return nil

}

func (repo *RepositoryStruct) loadKvartirasFromFile() error {
	// один мьютекс охватывает и файл и слайс
	repo.muKvartiras.Lock()
	defer repo.muKvartiras.Unlock()

	file, err := os.Open("kvartiras.csv")
	if err != nil {
		// если файла не существует значит первый запуск вероятнее всего, а значит дальнейшие действия по парсингу не имеют смысла
		if errors.Is(err, fs.ErrNotExist) { // os.IsNotExist старый метод, его не юзаем
			log.Println("файл kvartiras.csv не существует, первый запуск?")
			return nil
		}
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	count := 0

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		komnat, _ := strconv.Atoi(record[2])

		kvartira := &models.Kvartira{
			Id:     record[0],
			Number: record[1],
			Komnat: komnat,
		}

		repo.kvartiras = append(repo.kvartiras, kvartira)
		count++
	}

	log.Printf("%v квартир загружено из файла", count)

	return nil

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

func (repo *RepositoryStruct) UpdateFile(m models.ModelInterface) error {
	// важные моменты:
	// создавать временный файл, потом заменять им основной (в случае success). Иначе в случае ошибки закосячим основной файл
	// читаем файл по-строчно (на случай если вдруг большой файл будет)
	// после замененной строки не нужно перебирать строки, а делаем копирование оставшегося файла (так как id уникален в нашем случае)

	var searchId string
	var newLine []string
	var filename string
	var err error

	switch dataInt := m.(type) {
	case *models.Member:
		searchId = dataInt.Id
		newLine = []string{
			dataInt.Id,
			dataInt.Name,
			dataInt.Phone,
			strconv.Itoa(dataInt.Community),
		}

		filename = "members"
	case *models.Kvartira:
		searchId = dataInt.Id
		newLine = []string{
			dataInt.Id,
			dataInt.Number,
			strconv.Itoa(dataInt.Komnat),
		}

		filename = "kvartiras"
	default:
		return errors.New("неведомый тип данных")
	}

	file, err := os.Open(filename + ".csv")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) { // os.IsNotExist старый метод, его не юзаем
			return errors.New("файла с данными еще нет :(")
		}
		log.Printf("Файл есть, но чтение не удалось %v: %v", filename, err)
		return errors.New("файл с данными есть, но чтение не удалось :(")
	}
	defer file.Close()

	// читаем и записываем в tempfile чтоб в случае ошибки не сломать оригинальный файл
	tempfile, err := os.Create(filename + ".tmp")
	if err != nil {
		return errors.New("файл tmp не создан, выполнение операции невозможно")
	}
	defer tempfile.Close()

	//readfile := bufio.NewScanner(file)
	//writefile := bufio.NewWriter(tempfile)
	found := false

	readfile := csv.NewReader(file)
	writefile := csv.NewWriter(tempfile)

	for {
		record, err := readfile.Read()

		if err == io.EOF {
			// если конец файла
			break
		}

		if err != nil {
			return errors.New("ошибка чтения csv")
		}

		if len(record) > 0 {
			oldLineId := record[0]

			// если строка с искомым id есть, то заменяем ее на новую
			if oldLineId == searchId {
				record = newLine
				found = true
			}

			if err := writefile.Write(record); err != nil {
				return errors.New("запись измененной строки в файл не удалась")
			}

			// TODO: что-то не получилось записывать остаток файла, чтоб не перебирать строки бессмысленно. Нужно разобраться почему не срабатывает
			// // чтение остатка файла начнется с текущего оффсета
			// if _, err := io.Copy(writefile, file); err != nil {
			// 	return errors.New("не удалось записать остаток файла")
			// } else {
			// 	log.Println("записал остаток файла")
			// }
			// break // останавливаем for

		}

	}

	if !found {
		return errors.New("id не найден в файле")
	}

	// обязательно делаем Flush для гарантии полной очистки буфера!
	writefile.Flush()
	if err := writefile.Error(); err != nil {
		return errors.New("flush вернул ошибку")
	}

	// закрыты должны быть оба файла. defer'ы оставляем на случай ошибок (ошибки не будет при повторной попытке закрытия)
	if err := file.Close(); err != nil {
		return errors.New("ошибка закрытия основного файла")
	}
	if err := tempfile.Close(); err != nil {
		return errors.New("ошибка закрытия временного файла")
	}

	if err := os.Rename(filename+".tmp", filename+".csv"); err != nil {
		return errors.New("не удалось переименовать файл")
	}

	return nil

}

func (repo *RepositoryStruct) UpdateSlice(m models.ModelInterface) error {

	switch data := m.(type) {
	case *models.Member:
		//log.Println("in model member")
		repo.muMembers.Lock()
		defer repo.muMembers.Unlock()

		//log.Printf("data id: %v", data.Id)

		for i, m := range repo.members {
			//log.Printf("%v", m)
			if m.Id == data.Id {
				repo.members[i] = data
				//log.Println("FOUND")
				break // так как id уникальный дальше можем не перебирать слайс
			}
		}

		check, _ := json.Marshal(repo.members)
		log.Println(string(check))

	case *models.Kvartira:
		repo.muKvartiras.Lock()
		defer repo.muKvartiras.Unlock()

		for i, m := range repo.kvartiras {
			if m.Id == data.Id {
				repo.kvartiras[i] = data
				break // так как id уникальный дальше можем не перебирать слайс
			}
		}

		check, _ := json.Marshal(repo.kvartiras)
		log.Println(string(check))

	default:
		return errors.New("неведомый тип save")
	}

	return nil

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

	switch data := m.(type) {
	case *models.Member:

		filename := "members.csv"

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("ошибка открытия файла %v для записи: %v", filename, err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		memberCsv := []string{
			data.Id,
			data.Name,
			data.Phone,
			strconv.Itoa(data.Community),
		}
		writer.Write(memberCsv)

	case *models.Kvartira:

		filename := "kvartiras.csv"

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("ошибка открытия файла %v для записи: %v", filename, err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		kvartiraCsv := []string{
			data.Id,
			data.Number,
			strconv.Itoa(data.Komnat),
		}
		writer.Write(kvartiraCsv)

	default:
		return errors.New("неведомый тип savetofile")
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
