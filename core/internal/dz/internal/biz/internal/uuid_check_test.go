package biz_internal

import "testing"

func TestBizStruct_UuidCheck(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			name: "valid uuid",
			b:    &BizStruct{},
			args: args{
				id: "be5cb2a0-911c-451e-84af-5ff276583d28",
			},
			wantErr: false,
		},
		{
			name: "empty uuid",
			b:    &BizStruct{},
			args: args{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "invalid uuid",
			b:    &BizStruct{},
			args: args{
				id: "be5cb2a05ff276583d28",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UuidCheck(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.uuidCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
