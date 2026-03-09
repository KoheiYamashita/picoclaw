package i18n

import "testing"

func TestNormalizeLocale(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "en"},
		{"en", "en"},
		{"ja", "ja"},
		{"ja-JP", "ja"},
		{"en_US", "en"},
		{"EN", "en"},
		{"  ja  ", "ja"},
	}
	for _, tt := range tests {
		got := NormalizeLocale(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeLocale(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestT(t *testing.T) {
	// English
	got := T("en", "status.thinking")
	if got != "Thinking..." {
		t.Errorf("T(en, status.thinking) = %q", got)
	}

	// Japanese
	got = T("ja", "status.thinking")
	if got != "思考中..." {
		t.Errorf("T(ja, status.thinking) = %q", got)
	}

	// Fallback to en for unknown locale
	got = T("fr", "status.thinking")
	if got != "Thinking..." {
		t.Errorf("T(fr, status.thinking) = %q, want English fallback", got)
	}

	// Unknown key returns the key
	got = T("en", "nonexistent.key")
	if got != "nonexistent.key" {
		t.Errorf("T(en, nonexistent.key) = %q, want key itself", got)
	}
}

func TestTf(t *testing.T) {
	got := Tf("en", "status.searching_q", "golang")
	if got != "Searching... (golang)" {
		t.Errorf("Tf(en, status.searching_q, golang) = %q", got)
	}

	got = Tf("ja", "status.searching_q", "golang")
	if got != "検索中...（golang）" {
		t.Errorf("Tf(ja, status.searching_q, golang) = %q", got)
	}
}

func TestConfigLabels(t *testing.T) {
	// Japanese config label
	got := T("ja", "Model")
	if got != "モデル" {
		t.Errorf("T(ja, Model) = %q, want モデル", got)
	}

	// English config label falls back to key (since English labels match struct tags)
	got = T("en", "Model")
	if got != "Model" {
		t.Errorf("T(en, Model) = %q, want Model", got)
	}
}

func TestAgentMessages(t *testing.T) {
	got := T("en", "agent.context_window_warning")
	if got == "agent.context_window_warning" {
		t.Error("expected English warning message, got key itself")
	}

	got = T("ja", "agent.context_window_warning")
	if got == "agent.context_window_warning" {
		t.Error("expected Japanese warning message, got key itself")
	}
}
