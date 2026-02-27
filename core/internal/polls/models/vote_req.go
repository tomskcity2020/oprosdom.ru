package polls_models

import "errors"

type UnsafeVoteReq struct {
	PollID int    `json:"poll_id"`
	Vote   string `json:"vote"`
}

type ValidVoteReq struct {
	PollID int    `json:"poll_id"`
	Vote   string `json:"vote"`
	Jti    string `json:"jti"`
}

func (r *UnsafeVoteReq) Validate() (*ValidVoteReq, error) {
	if r.PollID <= 0 {
		return nil, errors.New("poll_id must be >0")
	}

	if r.Vote != "za" && r.Vote != "protiv" {
		return nil, errors.New("unexpected vote value")
	}

	validReq := &ValidVoteReq{
		PollID: r.PollID,
		Vote:   r.Vote,
	}

	return validReq, nil
}
