package shared_validate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ASCII only, short",
			input:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			expected: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
		{
			name:     "Contains non-ASCII (emoji)",
			input:    "Mozilla/5.0 😀",
			expected: "Mozilla/5.0 ",
		},
		{
			name:     "Contains control character (newline)",
			input:    "Mozilla/5.0\nChrome",
			expected: "Mozilla/5.0Chrome",
		},
		{
			name:     "Long string > 512 ASCII",
			input:    strings.Repeat("A", 600),
			expected: strings.Repeat("A", 512),
		},
		{
			name:     "Only non-ASCII characters (with spaces)",
			input:    "Привет мир 🌍",
			expected: "  ",
		},

		{
			name:     "Mix of ASCII and non-ASCII",
			input:    "Hello🌍World",
			expected: "HelloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UserAgentSanitize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
