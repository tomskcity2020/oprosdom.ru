package old_handlers

// import (
// 	"errors"
// 	"fmt"

// 	shared "oprosdom.ru/monolith/internal/shared"
// 	users_dto "oprosdom.ru/monolith/internal/users/dto"
// 	users_service "oprosdom.ru/monolith/internal/users/service"
// 	users_repo "oprosdom.ru/monolith/internal/users/service/repo"
// )

// type MemberSignUpHandler struct {
// 	Repo users_repo.RepositoryInterface
// }

// func (obj *MemberSignUpHandler) parse(params map[string]interface{}) (*users_dto.RequestSignUpDTO, error) {
// 	var paramsNotJson users_dto.RequestSignUpDTO
// 	err := shared.ParseParams(params, &paramsNotJson)
// 	if err != nil {
// 		return nil, errors.New("parsing MemberSignUp failed")
// 	}
// 	return &paramsNotJson, nil
// }

// func (obj *MemberSignUpHandler) basicValidation(dto users_dto.RequestSignUpDTO) error {
// 	if dto.Phone == "" {
// 		return errors.New("phone is empty")
// 	}
// 	if dto.City == "" {
// 		return errors.New("city is empty")
// 	}
// 	if dto.HouseId == "" {
// 		return errors.New("house_id is empty")
// 	}
// 	if dto.AptNumber == "" {
// 		return errors.New("apt number is empty")
// 	}
// 	if dto.Number == 0 {
// 		return errors.New("city is empty")
// 	}
// 	return nil
// }

// // возвращаем интерфейс{} чтоб соответствовать интерфейсу RpcHandler, если будем передавать тип responseDto то будет несоответствие
// func (obj *MemberSignUpHandler) Handle(params map[string]interface{}) (interface{}, error) {

// 	// парсим в dto
// 	var requestDto *users_dto.RequestSignUpDTO
// 	requestDto, err := obj.parse(params)
// 	if err != nil {
// 		return nil, errors.New("Handle MemberSignUp Parse return error")
// 	}
// 	fmt.Printf("Parsing result: %+v", requestDto)
// 	fmt.Println()

// 	// первичная валидация (еще будет по бизнес-логике валидация)
// 	err = obj.basicValidation(*requestDto)
// 	if err != nil {
// 		return nil, errors.New("Handle MemberSignUp Validation error")
// 	}
// 	fmt.Printf("Validation result: %+v", requestDto)
// 	fmt.Println()

// 	// вызываем сервисный слой соответствующий (а в нем вызываем бизнес-логику!)
// 	// из сервисного слоя получаем соответствующую структуру responseDTO
// 	userService := users_service.NewUserService(obj.Repo)
// 	responseDto, err := userService.MemberSignUpService(*requestDto)
// 	if err != nil {
// 		return nil, errors.New("UserService MemberSignUpService return error")
// 	}

// 	// responseDto уже указатель, поэтому брать адрес в памяти через & не нужно
// 	return responseDto, nil
// }
