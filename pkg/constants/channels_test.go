package constants

import (
	"testing"
)

func TestIsInternalChannel(t *testing.T) {
	tests := []struct {
		channel string
		want    bool
	}{
		{"cli", true},
		{"system", true},
		{"subagent", true},
		{"discord", false},
		{"telegram", false},
		{"slack", false},
		{"", false},
		{"CLI", false}, // case-sensitive
	}

	for _, tt := range tests {
		t.Run(tt.channel, func(t *testing.T) {
			got := IsInternalChannel(tt.channel)
			if got != tt.want {
				t.Errorf("IsInternalChannel(%q) = %v, want %v", tt.channel, got, tt.want)
			}
		})
	}
}
