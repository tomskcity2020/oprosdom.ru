package users_repo

type RepositoryInterface interface {
	Action() (*string, error)
}
