package elections

type Candidate struct {
	//Member
	role           string // house_manager / money_manager
	about          string
	program        string
	rating_za      int
	rating_unknown int
	rating_protiv  int
	status         string // candidate / on_vote / self_removed  / active  / revoked
	created_epoch  int64
}
