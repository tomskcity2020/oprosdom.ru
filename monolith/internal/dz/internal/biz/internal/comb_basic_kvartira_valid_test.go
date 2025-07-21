package biz_internal

import (
	"testing"

	"oprosdom.ru/monolith/internal/dz/internal/models"
)

func TestBizStruct_BasicKvartiraValidation(t *testing.T) {
	type args struct {
		kvartira *models.Kvartira
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
				kvartira: &models.Kvartira{
					Id:     "be5cb2a0-911c-451e-84af-5ff276583d28",
					Number: "52",
					Komnat: 5,
				},
			},
			wantErr: false,
		},
		{
			// если передадим заведомо невалидный id, а все остальные валидные данные, то в случае невозврата ошибки это будет означать то, что проверка на uuid закомментирована в BasicMemberValidation
			name: "test UuidCheck",
			b:    &BizStruct{},
			args: args{
				kvartira: &models.Kvartira{
					Id:     "",
					Number: "52",
					Komnat: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test kvNumberCheck",
			b:    &BizStruct{},
			args: args{
				kvartira: &models.Kvartira{
					Id:     "be5cb2a0-911c-451e-84af-5ff276583d28",
					Number: "",
					Komnat: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test kvKomnatCheck",
			b:    &BizStruct{},
			args: args{
				kvartira: &models.Kvartira{
					Id:     "be5cb2a0-911c-451e-84af-5ff276583d28",
					Number: "52",
					Komnat: -1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.BasicKvartiraValidation(tt.args.kvartira); (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.BasicKvartiraValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
