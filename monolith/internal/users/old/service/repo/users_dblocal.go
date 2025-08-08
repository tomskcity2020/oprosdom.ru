package users_repo

type Postgresql struct{}

func (db *Postgresql) Action() (*string, error) {
	data := "postgresql data"
	return &data, nil
}
