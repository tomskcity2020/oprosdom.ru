package users_repo

type Grpc struct{}

func (db *Grpc) Action() (*string, error) {
	data := "grpc data"
	return &data, nil
}
