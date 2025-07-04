package users_service

import (
	users_repo "oprosdom.ru/monolith/internal/users/service/repo"
)

type UserService struct {
	repo users_repo.RepositoryInterface
}

func NewUserService(repo users_repo.RepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}
