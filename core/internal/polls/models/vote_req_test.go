package polls_models

import (
	"errors"
	"testing"
)

func TestUnsafeVoteReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     UnsafeVoteReq
		want    *ValidVoteReq
		wantErr error
	}{
		{
			name: "valid za vote",
			req:  UnsafeVoteReq{PollID: 1, Vote: "za"},
			want: &ValidVoteReq{PollID: 1, Vote: "za"},
		},
		{
			name: "valid protiv vote",
			req:  UnsafeVoteReq{PollID: 2, Vote: "protiv"},
			want: &ValidVoteReq{PollID: 2, Vote: "protiv"},
		},
		{
			name:    "zero poll_id",
			req:     UnsafeVoteReq{PollID: 0, Vote: "za"},
			wantErr: errors.New("poll_id must be >0"),
		},
		{
			name:    "negative poll_id",
			req:     UnsafeVoteReq{PollID: -1, Vote: "protiv"},
			wantErr: errors.New("poll_id must be >0"),
		},
		{
			name:    "invalid vote value",
			req:     UnsafeVoteReq{PollID: 3, Vote: "abstain"},
			wantErr: errors.New("unexpected vote value"),
		},
		{
			name:    "empty vote value",
			req:     UnsafeVoteReq{PollID: 4, Vote: ""},
			wantErr: errors.New("unexpected vote value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.req.Validate()

			// Проверка ошибок
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка валидного результата
			if tt.wantErr == nil {
				if got.PollID != tt.want.PollID || got.Vote != tt.want.Vote {
					t.Errorf("Validate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
