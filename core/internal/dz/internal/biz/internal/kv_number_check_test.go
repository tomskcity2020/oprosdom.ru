package biz_internal

import "testing"

func TestBizStruct_kvNumberCheck(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			name: "valid kv number",
			b:    &BizStruct{},
			args: args{
				number: "125",
			},
			wantErr: false,
		},
		{
			name: "valid kv number 2",
			b:    &BizStruct{},
			args: args{
				number: "125a",
			},
			wantErr: false,
		},
		{
			name: "empty kv number",
			b:    &BizStruct{},
			args: args{
				number: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.kvNumberCheck(tt.args.number); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.kvNumberCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
