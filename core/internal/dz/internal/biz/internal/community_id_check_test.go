package biz_internal

import "testing"

func TestBizStruct_communityIdCheck(t *testing.T) {
	type args struct {
		communityId int
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			name: "valid community id",
			b:    &BizStruct{},
			args: args{
				communityId: 25,
			},
			wantErr: false,
		},
		{
			name: "not valid community id",
			b:    &BizStruct{},
			args: args{
				communityId: 0,
			},
			wantErr: true,
		},
		{
			name: "not valid community id 2",
			b:    &BizStruct{},
			args: args{
				communityId: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.communityIdCheck(tt.args.communityId); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.communityIdCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
