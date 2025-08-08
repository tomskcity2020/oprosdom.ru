package biz_internal

import "testing"

func TestBizStruct_nameCheck(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			name: "valid name",
			b:    &BizStruct{},
			args: args{
				name: "Сидоров",
			},
			wantErr: false,
		},
		{
			name: "valid name 2",
			b:    &BizStruct{},
			args: args{
				name: "Иванов Иван Иванович",
			},
			wantErr: false,
		},
		{
			name: "valid name 3",
			b:    &BizStruct{},
			args: args{
				name: "Fullname",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			b:    &BizStruct{},
			args: args{
				name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.nameCheck(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.nameCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
