package biz_internal

import (
	"testing"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func TestBizStruct_BasicMemberValidation(t *testing.T) {
	type args struct {
		member *models.Member
	}
	tests := []struct {
		name    string
		b       *BizStruct
		args    args
		wantErr bool
	}{
		{
			// здесь мы убеждаемся что валидные данные не вернут error. Более детальные тесты тут не имеют смысла так как для каждого внутреннего метода есть свой юнит тест
			name: "valid all",
			b:    &BizStruct{},
			args: args{
				member: &models.Member{
					Id: "be5cb2a0-911c-451e-84af-5ff276583d28",
					Name:      "Name Lastname",
					Phone:     "+79991231234",
					Community: 5,
				},
			},
			wantErr: false,
		},
		{
			// если передадим заведомо невалидный id, а все остальные валидные данные, то в случае невозврата ошибки это будет означать то, что проверка на uuid закомментирована в BasicMemberValidation
			name: "test UuidCheck",
			b:    &BizStruct{},
			args: args{
				member: &models.Member{
					Id: "",
					Name:      "Name Lastname",
					Phone:     "+79991231234",
					Community: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test nameCheck",
			b:    &BizStruct{},
			args: args{
				member: &models.Member{
					Id: "be5cb2a0-911c-451e-84af-5ff276583d28",
					Name:      "",
					Phone:     "+79991231234",
					Community: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test phoneCheck",
			b:    &BizStruct{},
			args: args{
				member: &models.Member{
					Id: "be5cb2a0-911c-451e-84af-5ff276583d28",
					Name:      "Name Lastname",
					Phone:     "",
					Community: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test communityIdCheck",
			b:    &BizStruct{},
			args: args{
				member: &models.Member{
					Id: "be5cb2a0-911c-451e-84af-5ff276583d28",
					Name:      "Name Lastname",
					Phone:     "+79991231234",
					Community: -1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.BasicMemberValidation(tt.args.member); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.BasicMemberValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}