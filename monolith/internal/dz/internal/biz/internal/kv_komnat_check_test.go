package biz_internal

import "testing"

func TestBizStruct_kvKomnatCheck(t *testing.T) {
	type args struct {
		komnat int
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			name: "valid kv komnat",
			b:    &BizStruct{},
			args: args{
				komnat: 3,
			},
			wantErr: false,
		},
		{
			name: "not valid kv komnat",
			b:    &BizStruct{},
			args: args{
				komnat: 0,
			},
			wantErr: true,
		},
		{
			name: "not valid kv komnat 2",
			b:    &BizStruct{},
			args: args{
				komnat: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.kvKomnatCheck(tt.args.komnat); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.kvKomnatCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
