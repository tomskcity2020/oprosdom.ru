package service_internal

import (
	"fmt"
	"log"
	"sync"
	"time"

	"oprosdom.ru/monolith/internal/dz/internal/biz"
	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
)

type ServiceStruct struct {
	repo repo.RepositoryInterface
	biz biz.BizInterface
}

func NewCallInternalService(repo repo.RepositoryInterface, biz biz.BizInterface) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
		biz: biz,
	}
}

func (s *ServiceStruct) CountData() {

	s.repo.LoadFromFile("members.json")
	s.repo.LoadFromFile("kvartiras.json")

	showMembers := len(s.repo.GetSliceMembers())
	fmt.Printf("Всего участников: %v", showMembers)
	fmt.Println()

	showKvartiras := len(s.repo.GetSliceKvartiras())
	fmt.Printf("Всего квартир: %v", showKvartiras)
	fmt.Println()
}

func (s *ServiceStruct) RunParallel(modelsData []models.ModelInterface) {

	// чтобы гарантировать корректную работу функции в методе CheckSlices нужно стартануть его ДО запуска кода где вставляются данные в слайс
	// иначе может быть не пустой слайс когда начнет исполняться CheckSlices и повлияет на корректный вывод изменений слайса

	// wgCheck := sync.WaitGroup{}
	// wgCheck.Add(1)
	// go func() {
	// 	wgCheck.Done() // defer тут не юзаем иначе CheckSlices будем ждать бесконечно
	// 	s.CheckSlices()
	// }()
	// wgCheck.Wait() // ждем когда горутина с CheckSlices стартанет

	ch := make(chan models.ModelInterface, 5)
	//ch := make(chan models.ModelInterface)

	// СТАРТУЕМ ЧИТАТЕЛЕЙ
	wgRead := sync.WaitGroup{}
	wgRead.Add(len(modelsData))
	//log.Printf("len modelsData %+v", len(modelsData))

	wgCheckGorutReadStarts := sync.WaitGroup{} // делаем войтгруппу для того, чтобы убедиться что горутина read стартанула раньше, чем начнется старт writer'ов
	wgCheckGorutReadStarts.Add(1)

	// если цикл чтения запущен до того как появился хотя бы один элемент в канале, он будет ожидать поступления данных, а не завершится сразу
	go func() { // запускаем в отдельной горутине чтоб не заблокировать основную горутину иначе словим deadlock так как, эта горутина переходит в состояние ожидания
		wgCheckGorutReadStarts.Done() // разрешаем wg до того как начинаем читать канал
		for m := range ch {
			// wgRead.Add(1) неправильно! нужно считать элементы до вызова горутины
			go func(model models.ModelInterface) { // при появлении в канале записи стартуем новую горутину потому, что save(model) может занимать какое-то время (в теории), чтоб выполнять параллельно множество операций,а не последовательно
				defer wgRead.Done()
				// нам нужно передать контекст непосредственно в Save чтоб там прервать операцию если вдруг долго исполняется
				s.repo.Save(model)

			}(m)

		}
	}()

	wgCheckGorutReadStarts.Wait()

	// СТАРТУЕМ ПИСАТЕЛЕЙ
	// размер буф канала не должен быть равен количеству записываемых элементов в него: если буф канал = 5, то это значит что это буфер на 5 элементов, чем больше буфер, тем меньше блокировок будет. Но слишком большой буфер это RAM, поэтому нужно по ситуации смотреть

	wgWrite := sync.WaitGroup{}
	//wgWrite.Add(len(modelsData))

	for i, model := range modelsData {

		wgWrite.Add(1)

		go func(wn int, m models.ModelInterface) {
			defer wgWrite.Done()
			//fmt.Printf("Писатель №%v начал запись в канал\n", wn)
			ch <- m
		}(i+1, model)

		// Println не юзаем потому, что это отдельная операция и при множестве горутин форматирование нарушается потому что они могут вклиниться перед Println
		//fmt.Println()

	}

	wgWrite.Wait()
	close(ch) // закрываем канал только после записи в него всех данных писателями

	wgRead.Wait() // ждем завершения работы читателей

	// теоретически у нас всегда будет >0, так как как минимум один экземпляр создастся интерактивно, но на всякий добавляем проверку
	if s.repo.MembersInSliceNow() > 0 {
		s.repo.SaveToFile("members.json")
	}

	if s.repo.KvartirasInSliceNow() > 0 {
		s.repo.SaveToFile("kvartiras.json")
	}

}

func (s *ServiceStruct) RunSeq(modelsData []models.ModelInterface) {
	wgCheck := sync.WaitGroup{}
	wgCheck.Add(1)
	go func() {
		wgCheck.Done() // defer тут не юзаем иначе CheckSlices будем ждать бесконечно
		s.CheckSlices()
	}()
	wgCheck.Wait() // ждем когда горутина с CheckSlices стартанет

	ch := make(chan models.ModelInterface)

	wgWrite := sync.WaitGroup{}
	wgWrite.Add(1)
	go func(md []models.ModelInterface) {
		defer wgWrite.Done()
		defer close(ch)
		for _, m := range md {
			ch <- m
		}
	}(modelsData)

	wgRead := sync.WaitGroup{}
	wgRead.Add(1)
	go func() {
		defer wgRead.Done()
		for m := range ch {
			s.repo.Save(m)
		}
	}()

	wgWrite.Wait()
	wgRead.Wait()
	// канал уже закрыт defer'ом в горутине записи

	showMembers := len(s.repo.GetSliceMembers())
	fmt.Printf("Всего участников: %v", showMembers)
	fmt.Println()

	showKvartiras := len(s.repo.GetSliceKvartiras())
	fmt.Printf("Всего квартир: %v", showKvartiras)
	fmt.Println()

}

func (s *ServiceStruct) CheckSlices() {
	//fmt.Println("check starts")
	// узнавать количество структур в каждом слайсе на старте (относительно там должен быть 0 - если эта горутина стартанет раньше других)
	prevMembersCount := len(s.repo.GetSliceMembers())
	prevKvartirasCount := len(s.repo.GetSliceKvartiras())

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop() // обязательно освобождаем ресурсы, сработает когда основная горутина завершится

	//log.Println("test1")

	for range ticker.C {

		// узнавать количество структур в каждом слайсе каждые 200 мс

		nowMembers := s.repo.GetSliceMembers()
		nowKvartiras := s.repo.GetSliceKvartiras()

		nowMembersCount := len(nowMembers)
		nowKvartirasCount := len(nowKvartiras)

		if nowMembersCount > prevMembersCount {
			newMembers := nowMembers[prevMembersCount:] // все новые это те, индекс которых идет следом за prevMembersCount
			for _, member := range newMembers {
				log.Printf("New value in Members: %+v ", member)
			}
			prevMembersCount = nowMembersCount
		}

		if nowKvartirasCount > prevKvartirasCount {
			newKvartiras := nowKvartiras[prevKvartirasCount:]
			for _, kvartira := range newKvartiras {
				log.Printf("New value in Kvartiras: %+v ", kvartira)
			}
			prevKvartirasCount = nowKvartirasCount
		}

	}
}
