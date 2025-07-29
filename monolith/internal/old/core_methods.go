package core

import (
	users_handlers "oprosdom.ru/monolith/internal/users/handlers"
	users_repo "oprosdom.ru/monolith/internal/users/service/repo"
)

type RpcHandler interface {
	Handle(map[string]interface{}) (interface{}, error)
}

var Handlers = make(map[string]RpcHandler)

func RegisterHandlers() {

	var usersRepo users_repo.RepositoryInterface
	usersRepo = &users_repo.Postgresql{}

	Handlers["member_signup"] = &users_handlers.MemberSignUpHandler{Repo: usersRepo}
	Handlers["change_phone"] = &users_handlers.ChangePhone{}
}

// инициализация обработчиков через init(). Если что можно вынести инициализацию в main
func init() {
	RegisterHandlers()
}
