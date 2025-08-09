package biz_internal

import (
	"testing"

	"github.com/google/uuid"
)

func TestBizStruct_UuidCreate(t *testing.T) {
	tests := []struct {
		name    string
		b       *BizStruct
		wantErr bool
	}{
		{
			name:    "success uuid generation",
			b:       &BizStruct{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.UuidCreate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BizStruct.UuidCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if _, err := uuid.Parse(got); err != nil {
					t.Errorf("UuidCreate returned invalid uuid: %v", got)
				}
			}
		})
	}
}
