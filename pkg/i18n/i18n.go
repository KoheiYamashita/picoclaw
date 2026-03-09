package i18n

import (
	"fmt"
	"strings"
)

// messages holds per-locale message catalogs.
// map[locale]map[key]translation
var messages = map[string]map[string]string{}

// register adds translations for a locale.
func register(locale string, m map[string]string) {
	if messages[locale] == nil {
		messages[locale] = map[string]string{}
	}
	for k, v := range m {
		messages[locale][k] = v
	}
}

// NormalizeLocale extracts a two-letter language code from a locale string.
// Examples: "ja-JP" -> "ja", "en_US" -> "en", "ja" -> "ja", "" -> "en".
func NormalizeLocale(locale string) string {
	locale = strings.TrimSpace(locale)
	if locale == "" {
		return "en"
	}
	// Handle both "ja-JP" and "ja_JP" forms
	for _, sep := range []string{"-", "_"} {
		if idx := strings.Index(locale, sep); idx > 0 {
			locale = locale[:idx]
		}
	}
	return strings.ToLower(locale)
}

// T returns the localized string for the given key and locale.
// Falls back to English ("en"), then returns the key itself.
func T(locale, key string) string {
	locale = NormalizeLocale(locale)

	if m, ok := messages[locale]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	// Fallback to English
	if locale != "en" {
		if m, ok := messages["en"]; ok {
			if v, ok := m[key]; ok {
				return v
			}
		}
	}
	return key
}

// Tf returns the localized string with fmt.Sprintf formatting.
func Tf(locale, key string, args ...interface{}) string {
	return fmt.Sprintf(T(locale, key), args...)
}
