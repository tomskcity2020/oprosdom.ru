package models

import (
	"testing"
)

func TestUnsafeCodeCheckReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     UnsafeCodeCheckReq
		want    *ValidatedCodeCheckReq
		wantErr bool
	}{
		{
			name: "valid request",
			req: UnsafeCodeCheckReq{
				Phone: "+79994567890",
				Code:  1234,
			},
			want: &ValidatedCodeCheckReq{
				Phone: "+79994567890",
				Code:  1234,
			},
			wantErr: false,
		},
		{
			name: "invalid phone number",
			req: UnsafeCodeCheckReq{
				Phone: "123",
				Code:  1234,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "code too short",
			req: UnsafeCodeCheckReq{
				Phone: "+79994567890",
				Code:  123,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "code too long",
			req: UnsafeCodeCheckReq{
				Phone: "+79994567890",
				Code:  12345,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty phone",
			req: UnsafeCodeCheckReq{
				Phone: "",
				Code:  1234,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UnsafeCodeCheckReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.Phone != tt.want.Phone {
				t.Errorf("UnsafeCodeCheckReq.Validate() phone = %v, want %v", got.Phone, tt.want.Phone)
			}
			if got.Code != tt.want.Code {
				t.Errorf("UnsafeCodeCheckReq.Validate() code = %v, want %v", got.Code, tt.want.Code)
			}
		})
	}
}
