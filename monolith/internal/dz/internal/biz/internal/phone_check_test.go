package biz_internal

import "testing"

func TestBizStruct_phoneCheck(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			name: "valid phone",
			b:    &BizStruct{},
			args: args{
				phone: "+7(999)123-1234",
			},
			wantErr: false,
		},
		{
			name: "valid phone 2",
			b:    &BizStruct{},
			args: args{
				phone: "89991231234",
			},
			wantErr: false,
		},
		{
			name: "valid phone 3",
			b:    &BizStruct{},
			args: args{
				phone: "7-999-1231234",
			},
			wantErr: false,
		},
		{
			name: "empty phone",
			b:    &BizStruct{},
			args: args{
				phone: "",
			},
			wantErr: true,
		},
		{
			name: "incorrect phone",
			b:    &BizStruct{},
			args: args{
				phone: "123",
			},
			wantErr: true,
		},
		{
			name: "incorrect phone 2",
			b:    &BizStruct{},
			args: args{
				phone: "+799912312345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.phoneCheck(tt.args.phone); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.phoneCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
