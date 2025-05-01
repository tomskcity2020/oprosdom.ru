package model

type Member_personal struct {
	id        int
	firstname string
	lastname  string
}

// сеттер
func (m *Member_personal) Set_Firstname(firstname string) {
	// здесь же можем сделать валидацию
	m.firstname = firstname
}

// геттер
func (m *Member_personal) Firstname() string {
	return m.firstname
}
