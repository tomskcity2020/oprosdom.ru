package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeMsg_Validate(t *testing.T) {
	tests := []struct {
		name        string
		input       UnsafeMsg
		expected    *ValidatedMsg
		expectedErr string
	}{
		{
			name: "valid input with retry 1",
			input: UnsafeMsg{
				Phone: "+79123456789",
				Code:  1234,
				Retry: 1,
			},
			expected: &ValidatedMsg{
				Phone: "+79123456789",
				Code:  1234,
				Retry: 1,
				Type:  "mobile",
			},
			expectedErr: "",
		},
		{
			name: "valid input with retry 3",
			input: UnsafeMsg{
				Phone: "+73822905090",
				Code:  9999,
				Retry: 3,
			},
			expected: &ValidatedMsg{
				Phone: "+73822905090",
				Code:  9999,
				Retry: 3,
				Type:  "landline",
			},
			expectedErr: "",
		},
		{
			name: "invalid phone number",
			input: UnsafeMsg{
				Phone: "invalid",
				Code:  1234,
				Retry: 1,
			},
			expected:    nil,
			expectedErr: "incorrect_format_phone", // Updated to match actual error
		},
		{
			name: "code too low",
			input: UnsafeMsg{
				Phone: "+79123456789",
				Code:  999,
				Retry: 1,
			},
			expected:    nil,
			expectedErr: "code must be 1000 to 9999",
		},
		{
			name: "code too high",
			input: UnsafeMsg{
				Phone: "+79123456789",
				Code:  10000,
				Retry: 1,
			},
			expected:    nil,
			expectedErr: "code must be 1000 to 9999",
		},
		{
			name: "retry zero - invalid",
			input: UnsafeMsg{
				Phone: "+79123456789",
				Code:  1234,
				Retry: 0,
			},
			expected:    nil,
			expectedErr: "retry must be 1 or 2 or 3",
		},
		{
			name: "retry too high",
			input: UnsafeMsg{
				Phone: "+79123456789",
				Code:  1234,
				Retry: 4,
			},
			expected:    nil,
			expectedErr: "retry must be 1 or 2 or 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.Validate()

			if tt.expectedErr != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, result) {
					assert.Equal(t, tt.expected, result)
				}
			}
		})
	}
}
