package shared_validate

import (
	"net"
	"testing"
)

func TestIpValidate(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		wantErr bool
	}{
		{
			name:    "valid IPv4",
			args:    args{p: "192.168.1.1"},
			want:    net.ParseIP("192.168.1.1"),
			wantErr: false,
		},
		{
			name:    "valid IPv6 short",
			args:    args{p: "2001:db8::1"},
			want:    net.ParseIP("2001:db8::1"),
			wantErr: false,
		},
		{
			name:    "valid IPv6 full (should compress internally)",
			args:    args{p: "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			want:    net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
			wantErr: false,
		},
		{
			name:    "invalid IP",
			args:    args{p: "not-an-ip"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty IP",
			args:    args{p: ""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "malformed IPv6",
			args:    args{p: "2001:::1"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IpValidate(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("IpValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil && got != nil {
				t.Errorf("IpValidate() = %v, want nil", got)
				return
			}

			if tt.want != nil && (got == nil || !got.Equal(tt.want)) {
				t.Errorf("IpValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}
