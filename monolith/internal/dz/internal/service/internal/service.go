package service_internal

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
)

type ServiceStruct struct {
	repo repo.RepositoryInterface
}

func NewCallInternalService(repo repo.RepositoryInterface) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
	}
}

// основной принцип сервисного пакета: берем данные извне или из репо, перекидываем в бизнес-логику, полученный результат сохраняем в репо или отдаем вовне

func (obj *ServiceStruct) RunParallel(ctx context.Context, modelsData []models.ModelInterface) {

	// чтобы гарантировать корректную работу функции в методе CheckSlices нужно стартануть его ДО запуска кода где вставляются данные в слайс
	// иначе может быть не пустой слайс когда начнет исполняться CheckSlices и повлияет на корректный вывод изменений слайса

	wgCheck := sync.WaitGroup{}
	wgCheck.Add(1)
	go func() {
		wgCheck.Done() // defer тут не юзаем иначе CheckSlices будем ждать бесконечно
		obj.CheckSlices()
	}()
	wgCheck.Wait() // ждем когда горутина с CheckSlices стартанет

	ch := make(chan models.ModelInterface, 5)

	// СТАРТУЕМ ЧИТАТЕЛЕЙ
	wgRead := sync.WaitGroup{}

	// если цикл чтения запущен до того как появился хотя бы один элемент в канале, он будет ожидать поступления данных, а не завершится сразу
	go func() { // запускаем в отдельной горутине чтоб не заблокировать основную горутину иначе словим deadlock так как, эта горутина переходит в состояние ожидания
		for m := range ch {

			wgRead.Add(1)

			go func(model models.ModelInterface) { // при появлении в канале записи стартуем новую горутину потому, что save(model) может занимать какое-то время (в теории), чтоб выполнять параллельно множество операций,а не последовательно
				defer wgRead.Done()
				// нам нужно передать контекст непосредственно в Save чтоб там прервать операцию если вдруг долго исполняется
				obj.repo.Save(model)

			}(m)

		}
	}()

	// СТАРТУЕМ ПИСАТЕЛЕЙ
	// размер буф канала не должен быть равен количеству записываемых элементов в него: если буф канал = 5, то это значит что это буфер на 5 элементов, чем больше буфер, тем меньше блокировок будет. Но слишком большой буфер это RAM, поэтому нужно по ситуации смотреть

	wgWrite := sync.WaitGroup{}

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

	showMembers := len(obj.repo.GetSliceMembers())
	fmt.Printf("Всего участников: %v", showMembers)
	fmt.Println()

	showKvartiras := len(obj.repo.GetSliceKvartiras())
	fmt.Printf("Всего квартир: %v", showKvartiras)
	fmt.Println()

}

func (obj *ServiceStruct) RunSeq(ctx context.Context, modelsData []models.ModelInterface) {
	wgCheck := sync.WaitGroup{}
	wgCheck.Add(1)
	go func() {
		wgCheck.Done() // defer тут не юзаем иначе CheckSlices будем ждать бесконечно
		obj.CheckSlices()
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
			obj.repo.Save(m)
		}
	}()

	wgWrite.Wait()
	wgRead.Wait()

	showMembers := len(obj.repo.GetSliceMembers())
	fmt.Printf("Всего участников: %v", showMembers)
	fmt.Println()

	showKvartiras := len(obj.repo.GetSliceKvartiras())
	fmt.Printf("Всего квартир: %v", showKvartiras)
	fmt.Println()

}

func (obj *ServiceStruct) CheckSlices() {
	//fmt.Println("check starts")
	// узнавать количество структур в каждом слайсе на старте (относительно там должен быть 0 - если эта горутина стартанет раньше других)
	prevMembersCount := len(obj.repo.GetSliceMembers())
	prevKvartirasCount := len(obj.repo.GetSliceKvartiras())

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop() // обязательно освобождаем ресурсы, сработает когда основная горутина завершится

	//log.Println("test1")

	for range ticker.C {

		// узнавать количество структур в каждом слайсе каждые 200 мс

		nowMembers := obj.repo.GetSliceMembers()
		nowKvartiras := obj.repo.GetSliceKvartiras()

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
