package old_users_service

// // в этом слое мы не реализуем бизнес-логику! Вызываем отсюда бизнес-логику и передаем туда данные из базы, а также на callback'е из бизнес-логики сохраняем данные в базу

// import (
// 	"errors"
// 	"fmt"

// 	users_business "oprosdom.ru/core/internal/users/business"
// 	users_dto "oprosdom.ru/core/internal/users/dto"
// )

// func (obj *UserService) MemberSignUpService(dto users_dto.RequestSignUpDTO) (*users_dto.ResponseSignUpDTO, error) {

// 	fmt.Println("UserService MemberSignUpService in action")

// 	businessLogic := users_business.NewUsersBusinessLogic()

// 	// посмотреть что сюда приходит можно в соответствующем обработчике

// 	// делаем запрос в базу данных на основе dto и получаем нужные для бизнес-логики данные
// 	data, err := obj.repo.Action()
// 	if err != nil {
// 		return nil, errors.New("Data receiving has been failed")
// 	}
// 	fmt.Println(data)

// 	// отправляем данные в бизнес-слой и получаем ответ
// 	result, err := businessLogic.SignupBizLogic()
// 	if err != nil {
// 		return nil, errors.New("BusinessLogic has been failed")
// 	}
// 	fmt.Println(result)

// 	// формируем структуру с типом responsedto и возвращаем
// 	serviceResult := users_dto.ResponseSignUpDTO{
// 		PublicId: 1827437837,
// 	}

// 	// берем из базы данных данные юзера и создаем business entity из типа данных в users_business
// 	// MemberPersonal := &users_business.MemberPersonal{
// 	// 	Id:              2,
// 	// 	Fake_id:         24,
// 	// 	Firstname:       "Имя",
// 	// 	Lastname:        "Фамилия",
// 	// 	Mobile:          "+71231231234",
// 	// 	Mobile_verified: false,
// 	// 	Nickname:        "",
// 	// 	Status:          "new",
// 	// }

// 	// // затем передаем эти данные в бизнес логику
// 	// done, err := MemberPersonal.SetStatus("active")

// 	// if err != nil {
// 	// 	return errors.New("users business returned an error")
// 	// }

// 	// // далее сохраняем в базу данных статус = done, проверяем на ошибку и если все норм, то возвращаем result true

// 	// fmt.Printf("Новый статус %v присвоен юзеру id=%v", done, MemberPersonal.Fake_id)

// 	return &serviceResult, nil
// }
