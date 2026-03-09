package i18n

import (
	"regexp"
	"testing"
)

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
		// Accept-Language header formats
		{"ja, en;q=0.9", "ja"},
		{"en-US,en;q=0.5", "en"},
		{"ja-JP, en-US;q=0.8, fr;q=0.5", "ja"},
		{"en;q=1.0", "en"},
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
	// Japanese config label (namespaced with "config." prefix)
	got := T("ja", "config.Model")
	if got != "モデル" {
		t.Errorf("T(ja, config.Model) = %q, want モデル", got)
	}

	// English config label returns the struct tag value
	got = T("en", "config.Model")
	if got != "Model" {
		t.Errorf("T(en, config.Model) = %q, want Model", got)
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

// TestFormatSpecifierConsistency verifies that en and ja translations have
// matching format specifiers (%s, %d, etc.) to prevent runtime panics in Tf().
func TestFormatSpecifierConsistency(t *testing.T) {
	re := regexp.MustCompile(`%[sdvfgqxobt]`)

	enMessages := messages["en"]
	jaMessages := messages["ja"]

	for key, enVal := range enMessages {
		jaVal, ok := jaMessages[key]
		if !ok {
			continue // ja doesn't have this key; fallback to en is fine
		}

		enSpecs := re.FindAllString(enVal, -1)
		jaSpecs := re.FindAllString(jaVal, -1)

		if len(enSpecs) != len(jaSpecs) {
			t.Errorf("format specifier count mismatch for key %q: en has %d (%v), ja has %d (%v)",
				key, len(enSpecs), enSpecs, len(jaSpecs), jaSpecs)
		}
	}
}
