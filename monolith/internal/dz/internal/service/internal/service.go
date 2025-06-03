package service_internal

import (
	"fmt"
	"sync"

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

	// сначала запускаем функцию проверки слайсов каждые 200 мс
	obj.repo.Check()

	// хз почему, нужно разобраться: когда создаю канал после models:=[]models, то ругается на ModelInterface (не видит пакет models). Поэтому задаем создаем новый тип тут, который используем ниже
	//ch := make(chan models.ModelInterface)
	type ModelChannel chan models.ModelInterface
	type modelsModelInterface models.ModelInterface // еще один костыль, потому что ниже []models.ModelInterface не видит пакет

	fmt.Println("Сервис запущен. Пожалуйста, ожидайте..")

	models := []models.ModelInterface{
		models.NewUserFactory("Namme", "+79991231234", 54),
		models.NewUserFactory("Bobby", "+71000033214", 79),
		models.NewKvartiraFactory("135b", 3),
		models.NewUserFactory("Marry", "+72223342343", 91),
		models.NewKvartiraFactory("179", 1),
		models.NewKvartiraFactory("11a", 2),
		models.NewUserFactory("Alex", "+71000032344", 51),
	}

	totalElements := len(models)

	ch := make(ModelChannel, totalElements)
	wg := sync.WaitGroup{}

	// создание моделей и запись в канал
	for i, model := range models {

		wg.Add(1)

		go func(rtn_num int) {
			defer wg.Done() // а внутри Save делаем свою waitgroup и отслеживаем done на том уровне
			fmt.Printf("стартанула %v горутина\n", rtn_num)
			//obj.repo.Save(model)
			ch <- model
		}(i + 1)

		//time.Sleep(2 * time.Second)
		fmt.Printf("Передано в горутины: %v структура из %v \n", i+1, totalElements)
		// Println не юзаем потому, что это отдельная операция и при множестве горутин форматирование нарушается
		//fmt.Println()
	}
	wg.Wait() // важно!! сначала дожидаемся записи в канал, только потом его закрываем. Если наоборот - будет паника
	close(ch)

	// ================== REPO SAVE()

	wgSave := sync.WaitGroup{}

	for m := range ch {

		wgSave.Add(1)

		// несмотря на то, что канал у нас закрыт (данные в нем изменяться не будут и также учитывая что тут <-readonly) все равно передаем m аргументом, на всякий случай для доп проверки, потому что так рекомендуют
		go func(model modelsModelInterface) {
			defer wgSave.Done()
			obj.repo.Save(model)
		}(m)

	}

	wg.Wait()

	showMembers := obj.repo.Show("member")
	fmt.Printf("Всего участников: %v", showMembers)
	fmt.Println()

	showKvartiras := obj.repo.Show("kvartira")
	fmt.Printf("Всего квартир: %v", showKvartiras)
	fmt.Println()

}
