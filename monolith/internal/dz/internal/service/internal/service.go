package service_internal

import (
	"context"
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

	// создаем дочерний контекст для write
	writeChannelCtx, writeChannelCancel := context.WithCancel(ctx)
	defer writeChannelCancel()

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
	//ch := make(chan models.ModelInterface)

	// СТАРТУЕМ ЧИТАТЕЛЕЙ
	wgRead := sync.WaitGroup{}

	readTotal := len(modelsData)
	wgRead.Add(readTotal)

	wgCheckGorutReadStarts := sync.WaitGroup{} // делаем войтгруппу для того, чтобы убедиться что горутина read стартанула раньше, чем начнется старт writer'ов
	wgCheckGorutReadStarts.Add(1)

	// если цикл чтения запущен до того как появился хотя бы один элемент в канале, он будет ожидать поступления данных, а не завершится сразу
	go func(ctxGo context.Context, rt int) { // запускаем в отдельной горутине чтоб не заблокировать основную горутину иначе словим deadlock так как, эта горутина переходит в состояние ожидания
		// ВАЖНО!!!
		// контекст чтения из канала зависит от контекста записи в канал
		// контекст сохранения зависит от контекста чтения из канала
		// то есть при сигнале отмены сначала прекращаем запись в канал, затем прекращаем считывание из канала и после этого прекращаем сохранение

		readGoCtx, readGoCtxCancel := context.WithCancel(writeChannelCtx)
		defer readGoCtxCancel()

		saveGoCtx, saveGoCtxCancel := context.WithCancel(readGoCtx)
		defer saveGoCtxCancel()

		wgCheckGorutReadStarts.Done() // разрешаем wg до того как начинаем читать канал

		// вводим счетчик для того, чтобы в случае readGoCtx.Done() доделать оставшиеся wgRead.Done() иначе повиснем
		remaining := rt
		defer func() {
			for i := 0; i < remaining; i++ {
				log.Println("делаю wgRead.Done после отмены")
				wgRead.Done()
			}
		}()

		// эта waitgroup создана специально для того, чтобы defer'ы объявленные в начале этой горутины исполнялись только после wg.Wait иначе они преждевременно отменят контекст в репо (так как он дочерний)
		saveWg := sync.WaitGroup{}

		for m := range ch {
			select {
			case <-readGoCtx.Done():
				log.Println("прерываю read из канала")
			default:
				// для теста прерывания write:
				time.Sleep(500 * time.Millisecond)
				log.Println("считал значение из канала и передаю в репо")

				remaining--
				saveWg.Add(1)
				// wgRead.Add(1) неправильно! нужно считать элементы до вызова горутины
				go func(saveCtx context.Context, model models.ModelInterface) { // при появлении в канале записи стартуем новую горутину потому, что save(model) может занимать какое-то время (в теории), чтоб выполнять параллельно множество операций,а не последовательно
					defer wgRead.Done()
					// нам нужно передать контекст непосредственно в Save чтоб там прервать операцию если вдруг долго исполняется
					obj.repo.Save(saveCtx, model)

				}(saveGoCtx, m)
			}

		}
		saveWg.Wait()
	}(writeChannelCtx, readTotal)

	wgCheckGorutReadStarts.Wait()

	// СТАРТУЕМ ПИСАТЕЛЕЙ
	// размер буф канала не должен быть равен количеству записываемых элементов в него: если буф канал = 5, то это значит что это буфер на 5 элементов, чем больше буфер, тем меньше блокировок будет. Но слишком большой буфер это RAM, поэтому нужно по ситуации смотреть

	wgWrite := sync.WaitGroup{}
	//wgWrite.Add(len(modelsData))

	for i, model := range modelsData {
		select {
		case <-writeChannelCtx.Done():
			log.Println("GRACEFUL SHUTDOWN starts while RunParallel()")
			return
		default:
			// для теста прерывания write:
			time.Sleep(500 * time.Millisecond)

			wgWrite.Add(1)

			go func(wn int, m models.ModelInterface) {
				defer wgWrite.Done()
				log.Printf("Писатель №%v начал запись в канал\n", wn)
				ch <- m
			}(i+1, model)

		}
	}

	wgWrite.Wait()
	close(ch) // закрываем канал только после записи в него всех данных писателями

	wgRead.Wait() // ждем завершения работы читателей

	showMembers := len(obj.repo.GetSliceMembers())
	log.Printf("Всего участников: %v", showMembers)
	//fmt.Println()

	showKvartiras := len(obj.repo.GetSliceKvartiras())
	log.Printf("Всего квартир: %v", showKvartiras)
	//fmt.Println()

}

// func (obj *ServiceStruct) RunSeq(modelsData []models.ModelInterface) {
// 	wgCheck := sync.WaitGroup{}
// 	wgCheck.Add(1)
// 	go func() {
// 		wgCheck.Done() // defer тут не юзаем иначе CheckSlices будем ждать бесконечно
// 		obj.CheckSlices()
// 	}()
// 	wgCheck.Wait() // ждем когда горутина с CheckSlices стартанет

// 	ch := make(chan models.ModelInterface)

// 	wgWrite := sync.WaitGroup{}
// 	wgWrite.Add(1)
// 	go func(md []models.ModelInterface) {
// 		defer wgWrite.Done()
// 		defer close(ch)
// 		for _, m := range md {
// 			ch <- m
// 		}
// 	}(modelsData)

// 	wgRead := sync.WaitGroup{}
// 	wgRead.Add(1)
// 	go func() {
// 		defer wgRead.Done()
// 		for m := range ch {
// 			obj.repo.Save(m)
// 		}
// 	}()

// 	wgWrite.Wait()
// 	wgRead.Wait()
// 	// канал уже закрыт defer'ом в горутине записи

// 	showMembers := len(obj.repo.GetSliceMembers())
// 	fmt.Printf("Всего участников: %v", showMembers)
// 	fmt.Println()

// 	showKvartiras := len(obj.repo.GetSliceKvartiras())
// 	fmt.Printf("Всего квартир: %v", showKvartiras)
// 	fmt.Println()

// }

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
