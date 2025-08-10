package models

import (
	"net"
	"testing"
)

func TestUnsafeCodeCheckReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   UnsafeCodeCheckReq
		want    *ValidatedCodeCheckReq
		wantErr bool
	}{
		{
			name: "valid request with all fields",
			input: UnsafeCodeCheckReq{
				Phone:     "+79123456789",
				Code:      1234,
				UserAgent: "Mozilla/5.0",
				Ip:        "192.168.1.1",
			},
			want: &ValidatedCodeCheckReq{
				Phone:     "+79123456789",
				Code:      1234,
				UserAgent: "Mozilla/5.0",
				IP:        net.ParseIP("192.168.1.1"),
			},
			wantErr: false,
		},
		{
			name: "valid request with empty UserAgent",
			input: UnsafeCodeCheckReq{
				Phone:     "+79123456789",
				Code:      1234,
				UserAgent: "",
				Ip:        "192.168.1.1",
			},
			want: &ValidatedCodeCheckReq{
				Phone:     "+79123456789",
				Code:      1234,
				UserAgent: "",
				IP:        net.ParseIP("192.168.1.1"),
			},
			wantErr: false,
		},
		{
			name: "invalid phone number",
			input: UnsafeCodeCheckReq{
				Phone:     "invalid",
				Code:      1234,
				UserAgent: "Mozilla/5.0",
				Ip:        "192.168.1.1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "code too small",
			input: UnsafeCodeCheckReq{
				Phone:     "+79123456789",
				Code:      999,
				UserAgent: "",
				Ip:        "192.168.1.1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "code too big",
			input: UnsafeCodeCheckReq{
				Phone:     "+79123456789",
				Code:      10000,
				UserAgent: "Mozilla/5.0",
				Ip:        "192.168.1.1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid IP address",
			input: UnsafeCodeCheckReq{
				Phone:     "+79123456789",
				Code:      1234,
				UserAgent: "",
				Ip:        "invalid.ip",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.Phone != tt.want.Phone {
				t.Errorf("Phone = %v, want %v", got.Phone, tt.want.Phone)
			}

			if got.Code != tt.want.Code {
				t.Errorf("Code = %v, want %v", got.Code, tt.want.Code)
			}

			if got.UserAgent != tt.want.UserAgent {
				t.Errorf("UserAgent = %v, want %v", got.UserAgent, tt.want.UserAgent)
			}

			if !got.IP.Equal(tt.want.IP) {
				t.Errorf("IP = %v, want %v", got.IP, tt.want.IP)
			}
		})
	}
}
