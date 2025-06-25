package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/service"
)

func main() {

	// modelsData := []models.ModelInterface{
	// 	models.NewUserFactory("Namme", "+79991231234", 54),
	// 	models.NewUserFactory("Bobby", "+71000033214", 79),
	// 	models.NewKvartiraFactory("135b", 3),
	// 	models.NewUserFactory("Marry", "+72223342343", 91),
	// 	models.NewKvartiraFactory("179", 1),
	// 	models.NewKvartiraFactory("11a", 2),
	// 	models.NewUserFactory("Alex", "+71000032344", 51),
	// }

	// логика graceful shutdown у нас такая: дожидаемся записи в слайсы и прекращаем дальнейшее выполнение цикла, выходим. Потому что посреди записи в слайс нелогично делать graceful shutdown, потому что он таковым являться не будет ввиду того, что часть данных запишется в слайс, а часть возможно нет
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	input := bufio.NewScanner(os.Stdin)

	for {

		// каждую итерацию создаем новый modelsData в который записываем данные из файла и интерактивного ввода, иначе будет дублирование если вне for разместить
		modelsData := make([]models.ModelInterface, 0)

		serviceEntity := service.NewServiceFactory()

		// каждую итерацию показываем обновленные данные (читаем файл заново каждую итерацию)
		serviceEntity.CountData()

		fmt.Println("---Выберите действие---")
		fmt.Println("1. Добавить жителя")
		fmt.Println("2. Добавить квартиру")
		fmt.Println("3. Выйти")

		input.Scan()

		selected := strings.TrimSpace(input.Text())

		select {
		case <-ctx.Done():
			fmt.Println("Graceful Shutdown starts")
			return
		default:

			switch selected {
			case "1":
				member, err := createMemberInput(input)
				if err != nil {
					fmt.Printf("Ошибка %v\n", err)
					continue
				}
				modelsData = append(modelsData, member)
				serviceEntity.RunParallel(modelsData)
				//serviceEntity.RunSeq(modelsData)
			case "2":
				kvartira, err := createKvartiraInput(input)
				if err != nil {
					fmt.Printf("Ошибка %v\n", err)
					continue
				}
				modelsData = append(modelsData, kvartira)
				serviceEntity.RunParallel(modelsData)
				//serviceEntity.RunSeq(modelsData)
			case "3":
				return
			default:
				fmt.Println("Неверный выбор. Выберите цифру соответствующую требуемому действию")
			}

		}

	}

}

func createMemberInput(input *bufio.Scanner) (models.ModelInterface, error) {
	fmt.Println("---Добавление жителя---")

	fmt.Println("Введите имя:")
	input.Scan()
	name := strings.TrimSpace(input.Text())

	fmt.Println("Телефон:")
	input.Scan()
	phone := strings.TrimSpace(input.Text())

	fmt.Println("Номер сообщества:")
	input.Scan()
	communityString := strings.TrimSpace(input.Text())
	community, err := strconv.Atoi(communityString)
	if err != nil {
		return nil, fmt.Errorf("номер сообщества должен быть числом")
	}

	return models.NewUserFactory(name, phone, community), nil
}

func createKvartiraInput(input *bufio.Scanner) (models.ModelInterface, error) {
	fmt.Println("---Добавление квартиры---")

	fmt.Println("Введите номер квартиры:")
	input.Scan()
	number := strings.TrimSpace(input.Text())

	fmt.Println("Кол-во комнат:")
	input.Scan()
	roomsString := strings.TrimSpace(input.Text())
	rooms, err := strconv.Atoi(roomsString)
	if err != nil {
		return nil, fmt.Errorf("число комнат должно быть числом")
	}

	return models.NewKvartiraFactory(number, rooms), nil
}
