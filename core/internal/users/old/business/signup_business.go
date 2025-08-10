package users_business

// type MemberPersInterface interface {
// 	SetStatus(new_status string) (setted_status string, err error)
// }

func (obj *UsersBusinessLogic) SignupBizLogic() (*string, error) {
	result := "business result to save in db"
	return &result, nil
}

// type MemberPersonal struct {
// 	Id              int
// 	Fake_id         int
// 	Firstname       string
// 	Lastname        string
// 	Mobile          string
// 	Mobile_verified bool
// 	Nickname        string
// 	Status          string // new active limited spamban fake - если чела в одном из сообществ забанили за спам, то нужно spamban статус применять к базовой структуре, чтоб во всех сообществах накладывать ограничения
// }

// func (m *MemberPersonal) SetStatus(new_status string) (setted_status string, err error) {
// 	if new_status == "active" && m.Mobile_verified {
// 		setted_status := "active"
// 		return setted_status, nil
// 	} else {
// 		return "", err
// 	}
// }

// type Community struct {
// 	id            int
// 	houseguid     string // (фиас uniq)
// 	country       string
// 	region        string
// 	naspunkt      string
// 	address       string
// 	latitude      float64
// 	longitude     float64
// 	started_epoch int64
// 	users_count   int
// 	status        string // waiting5 / pending30 / success
// }

// // методы:
// // ShowCommunity
// // AddCommunity

// type House_manager struct {
// 	//Member
// 	Role_started_at int64
// }

// type Money_manager struct {
// 	Member
// 	Role_started_at int64
// 	Fssp_checked    bool
// }

// // // сеттер
// // func (m *Member_personal) SetFirstname(firstname string) {
// // 	// здесь же можем сделать валидацию
// // 	m.firstname = firstname
// // }

// // // геттер
// // func (m *Member_personal) Firstname() string {
// // 	return m.firstname
// // }

// // func (m *Member_personal) SetLastname(lastname string) {
// // 	m.lastname = lastname
// // }

// // func (m *Member_personal) Lastname() string {
// // 	return m.lastname
// // }

// // // для теста сеттер, потом удалить!
// // func (m *Member_personal) SetID(id int) {
// // 	m.id = id
// // }

// // func (m *Member_personal) ID() int {
// // 	return m.id
// // }

// // участник может состоять в нескольких домах, но так как верификация везде своя (старшим или сообществом), то разделяем персональные данные и данные по Жителю

// type Member struct {
// 	Member_personal
// 	community_id             int
// 	kv                       string // потому что может быть с буквой
// 	verified_by_card         bool
// 	verified_by_housemanager bool
// 	verified_by_community    bool
// 	temporary_ban            bool // чел может быть зареган в нескольких Сообществах, в одном что-то натворить и схлопотать бан, это значит что бан распространяется только на this сообщество
// }

// // сеттер
// func (m *Member) SetCommunityID(id int) {
// 	m.community_id = id
// }

// // геттер
// func (m *Member) CommunityID() int {
// 	return m.community_id
// }

// type Violation struct{}

// //методы:
// //Temporary_ban (Старший по дому)
// //Ban_member (Сообщество по результатам голосования)

// type CreateMemberDTO struct {
// 	Firstname string
// 	Lastname  string
// 	Community int
// }

// // метод для первичной проверки данных
// func (dto *CreateMemberDTO) Validate() error {
// 	if dto.Firstname == "" {
// 		return errors.New("firstname is required")
// 	}
// 	if dto.Lastname == "" {
// 		return errors.New("lastname is required")
// 	}
// 	if dto.Community == 0 {
// 		return errors.New("community is required")
// 	}
// 	return nil
// }

// type UpdateMemberDto struct {
// 	Id        int
// 	Firstname string
// 	Lastname  string
// }
