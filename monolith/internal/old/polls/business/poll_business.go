package poll

type Poll struct {
	id                      int
	created_by_user_id      int
	created_by_housemanager bool
	title                   string
	description             string
	za                      int
	protiv                  int
	plus                    int
	minus                   int
}

//методы:
// Show_polls
// Add_vote
// Add_plus
// Add_minus
