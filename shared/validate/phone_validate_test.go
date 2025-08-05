package shared_validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneValidate(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedPhone string
		expectedType  string
		expectErr     string
	}{
		{
			name:          "Valid Russian mobile",
			input:         "+7 (999) 123-45-67",
			expectedPhone: "+79991234567",
			expectedType:  "mobile",
		},
		{
			name:          "Valid megafon mobile",
			input:         "89252341234",
			expectedPhone: "+79252341234",
			expectedType:  "mobile",
		},
		{
			name:          "Valid tele2 mobile",
			input:         "89772341234",
			expectedPhone: "+79772341234",
			expectedType:  "mobile",
		},
		{
			name:          "Valid Russian landline",
			input:         "+7 (495) 123-45-67",
			expectedPhone: "+74951234567",
			expectedType:  "landline",
		},
		{
			name:      "Invalid: Kazakhstan number (starts with +7 but not RU)",
			input:     "+7 701 123 4567",
			expectErr: "not_valid_ru_phone_number",
		},
		{
			name:      "Invalid: too short",
			input:     "123",
			expectErr: "not_valid_ru_phone_number",
		},
		{
			name:      "Invalid: non-numeric characters",
			input:     "abcdef",
			expectErr: "incorrect_format_phone",
		},
		{
			name:      "Empty input",
			input:     "   ",
			expectErr: "empty_phone",
		},
		{
			name:          "Valid Russian mobile without plus",
			input:         "8 (926) 000-00-00",
			expectedPhone: "+79260000000",
			expectedType:  "mobile",
		},
		{
			name:          "Valid Russian landline short format",
			input:         "4951234567",
			expectedPhone: "+74951234567",
			expectedType:  "landline",
		},
		{
			name:          "Valid Russian number with extra spaces",
			input:         "   +7 495 123 45 67  ",
			expectedPhone: "+74951234567",
			expectedType:  "landline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, typ, err := PhoneValidate(tt.input)

			if tt.expectErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectErr)
				assert.Empty(t, formatted)
				assert.Empty(t, typ)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPhone, formatted)
				assert.Equal(t, tt.expectedType, typ)
			}
		})
	}
}
