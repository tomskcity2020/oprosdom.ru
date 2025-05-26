package dz6

import (
	"errors"
	"math/rand"
	"time"
)

func CreateModel(model string) (interface{}, error) {

	memberNames := []string{"Николай", "Алексей", "Александр", "Василий"}
	memberPhones := []string{"+79991231234", "+79992342345", "+79139992838", "+79068177474"}
	kvartiraNumbers := []string{"18", "124", "22a", "137б"}

	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	randomNumMb := generator.Intn(len(memberNames))
	randomNumMp := generator.Intn(len(memberPhones))
	randomNumKn := generator.Intn(len(kvartiraNumbers))

	var readyModel ToSliceInterface

	switch model {
	case "member":
		readyModel = &Member{
			Name:      memberNames[randomNumMb],
			Phone:     memberPhones[randomNumMp],
			Community: randomNumMp,
		}
	case "kvartira":
		readyModel = &Kvartira{
			Number: kvartiraNumbers[randomNumKn],
			Komnat: randomNumKn,
		}
	default:
		return nil, errors.New("no model found")
	}

	// передаем во вторую функцию
	result, err := ToSlice(readyModel)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return result, nil

}
