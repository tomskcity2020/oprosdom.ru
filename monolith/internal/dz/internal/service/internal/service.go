package service_internal

import (
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

// основной принцип сервисного пакета: берем данные из репо, перекидываем в бизнес-логику, полученный результат сохраняем в репо.
// TEST проверить удобно ли будет работать с горутинами в сервисном слое, чтоб не загромождать бизнес-слой, или нужно делать доп уровень абстракции

func (obj *ServiceStruct) Run() {

	// сначала запускаем в отдельной горутине функцию проверки слайсов каждые 200 мс
	// TODO нужно вынести в пакет businesslogic бизнес-логику, здесь оставить только код в соответствии с принципом сервисного пакета
	go func() {
		// узнавать количество структур в каждом слайсе на старте (однозначно там по 0)
		prevMembersCount := len(obj.repo.GetSliceMembers())
		prevKvartirasCount := len(obj.repo.GetSliceKvartiras())

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop() // обязательно освобождаем ресурсы, сработает когда основная горутина завершится

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
	}()

	//fmt.Println("Сервис запущен. Пожалуйста, ожидайте..")

	modelsData := []models.ModelInterface{
		models.NewUserFactory("Namme", "+79991231234", 54),
		models.NewUserFactory("Bobby", "+71000033214", 79),
		models.NewKvartiraFactory("135b", 3),
		models.NewUserFactory("Marry", "+72223342343", 91),
		models.NewKvartiraFactory("179", 1),
		models.NewKvartiraFactory("11a", 2),
		models.NewUserFactory("Alex", "+71000032344", 51),
	}

	totalElements := len(modelsData)

	ch := make(chan models.ModelInterface, totalElements)
	wg := sync.WaitGroup{}

	// ==================== 1. Запись структур в канал
	for i, model := range modelsData {

		wg.Add(1)

		go func(rtn_num int, modelGoRtn models.ModelInterface) {
			defer wg.Done() // а внутри Save делаем свою waitgroup и отслеживаем done на том уровне
			fmt.Printf("стартанула %v горутина\n", rtn_num)
			ch <- modelGoRtn // передаем model через аргумент на
		}(i+1, model)

		fmt.Printf("Передано в горутины: %v структура из %v \n", i+1, totalElements)
		// Println не юзаем потому, что это отдельная операция и при множестве горутин форматирование нарушается
		//fmt.Println()

	}
	wg.Wait() // важно!! сначала дожидаемся записи в канал, только потом его закрываем. Если наоборот - будет паника
	close(ch)

	// ================== 2. REPO SAVE()

	wgSave := sync.WaitGroup{}

	for m := range ch {

		wgSave.Add(1)

		// несмотря на то, что канал у нас закрыт (данные в нем изменяться не будут и также учитывая что тут <-readonly) все равно передаем m аргументом, на всякий случай для доп проверки, потому что так рекомендуют
		go func(model models.ModelInterface) {
			defer wgSave.Done()
			obj.repo.Save(model)
		}(m)

	}

	wgSave.Wait()

	// ================== END

	showMembers := len(obj.repo.GetSliceMembers())
	fmt.Printf("Всего участников: %v", showMembers)
	fmt.Println()

	showKvartiras := len(obj.repo.GetSliceKvartiras())
	fmt.Printf("Всего квартир: %v", showKvartiras)
	fmt.Println()

}
