package utils

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "short string unchanged",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "exact length unchanged",
			input:  "hello",
			maxLen: 5,
			want:   "hello",
		},
		{
			name:   "truncated with ellipsis",
			input:  "hello world!",
			maxLen: 8,
			want:   "hello...",
		},
		{
			name:   "maxLen 0",
			input:  "hello",
			maxLen: 0,
			want:   "",
		},
		{
			name:   "maxLen 1",
			input:  "hello",
			maxLen: 1,
			want:   "h",
		},
		{
			name:   "maxLen 3",
			input:  "hello",
			maxLen: 3,
			want:   "hel",
		},
		{
			name:   "maxLen 4 with ellipsis",
			input:  "hello world",
			maxLen: 4,
			want:   "h...",
		},
		{
			name:   "unicode characters",
			input:  "こんにちは世界",
			maxLen: 5,
			want:   "こん...",
		},
		{
			name:   "empty string",
			input:  "",
			maxLen: 10,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}
