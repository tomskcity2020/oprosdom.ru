package biz_internal

import "testing"

func TestBizStruct_DecimalCheck(t *testing.T) {
	type args struct {
		amount string
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			name: "valid 1",
			b:    &BizStruct{},
			args: args{
				amount: "123.45",
			},
			wantErr: false,
		},
		{
			name: "valid 2",
			b:    &BizStruct{},
			args: args{
				amount: "100",
			},
			wantErr: false,
		},
		{
			name: "valid 3",
			b:    &BizStruct{},
			args: args{
				amount: "0.00",
			},
			wantErr: false,
		},
		{
			name: "empty",
			b:    &BizStruct{},
			args: args{
				amount: "",
			},
			wantErr: true,
		},
		{
			name: "negative",
			b:    &BizStruct{},
			args: args{
				amount: "-50",
			},
			wantErr: true,
		},
		{
			name: "not valid 1",
			b:    &BizStruct{},
			args: args{
				amount: "100.123",
			},
			wantErr: true,
		},
		{
			name: "not valid 2",
			b:    &BizStruct{},
			args: args{
				amount: "asd.qwe",
			},
			wantErr: true,
		},
		{
			name: "not valid 3",
			b:    &BizStruct{},
			args: args{
				amount: "12.12.12",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.DecimalCheck(tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.DecimalCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
