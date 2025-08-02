package models

import (
	"testing"
)

// тест основан на том, что мы отправляем заведомо неправильные данные и ожидаем ошибки, а не nil. Если возвращается nil, то можно сделать вывод о том, что функция не вызывается
// с useragent ситуация немного другая ввиду того, что "" строка проходит проверку, мы допускаем это. Поэтому с ua наоборот отправляем валидный ua и если возвращается "", то можно сделать вывод о том что что-то не то с функцией санитизации ua

func TestUnsafePhoneSendReq_Validate(t *testing.T) {
	tests := []struct {
		name            string
		req             UnsafePhoneSendReq
		wantErr         bool
		expectUserAgent string // в valid all отправляем валидный useragent и если придет "" значит функция закомментирована или еще что, но нам главное что этого достаточно для дальнейшего разбирательства что произошло
	}{
		{
			name: "valid all",
			req: UnsafePhoneSendReq{
				Phone:     "+79123456789",
				UserAgent: "Mozilla/5.0",
				Ip:        "192.168.1.1",
			},
			wantErr:         false,
			expectUserAgent: "Mozilla/5.0",
		},
		{
			name: "invalid phone",
			req: UnsafePhoneSendReq{
				Phone:     "123",
				UserAgent: "Mozilla/5.0",
				Ip:        "192.168.1.1",
			},
			wantErr: true,
		},
		{
			name: "invalid ip",
			req: UnsafePhoneSendReq{
				Phone:     "+79123456789",
				UserAgent: "Mozilla/5.0",
				Ip:        "invalid.ip",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.req.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Если ошибки не ожидалось, проверяем обработку UserAgent
			if !tt.wantErr {
				if result.UserAgent != tt.expectUserAgent {
					t.Errorf("UserAgent not processed correctly: got %q, want %q",
						result.UserAgent, tt.expectUserAgent)
				}
			}
		})
	}
}
