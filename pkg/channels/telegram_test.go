package channels

import (
	"strings"
	"testing"
)

// --- escapeHTML ---

func TestEscapeHTML_Ampersand(t *testing.T) {
	got := escapeHTML("a & b")
	if got != "a &amp; b" {
		t.Errorf("got %q", got)
	}
}

func TestEscapeHTML_LessThan(t *testing.T) {
	got := escapeHTML("a < b")
	if got != "a &lt; b" {
		t.Errorf("got %q", got)
	}
}

func TestEscapeHTML_GreaterThan(t *testing.T) {
	got := escapeHTML("a > b")
	if got != "a &gt; b" {
		t.Errorf("got %q", got)
	}
}

func TestEscapeHTML_Empty(t *testing.T) {
	got := escapeHTML("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestEscapeHTML_Mixed(t *testing.T) {
	got := escapeHTML("<b>bold & italic</b>")
	if !strings.Contains(got, "&lt;") || !strings.Contains(got, "&gt;") || !strings.Contains(got, "&amp;") {
		t.Errorf("not all chars escaped: %q", got)
	}
}

func TestEscapeHTML_NoSpecialChars(t *testing.T) {
	input := "hello world 123"
	got := escapeHTML(input)
	if got != input {
		t.Errorf("expected no change, got %q", got)
	}
}

// --- extractCodeBlocks ---

func TestExtractCodeBlocks_None(t *testing.T) {
	result := extractCodeBlocks("no code blocks here")
	if len(result.codes) != 0 {
		t.Errorf("expected 0 code blocks, got %d", len(result.codes))
	}
	if result.text != "no code blocks here" {
		t.Errorf("text should be unchanged, got %q", result.text)
	}
}

func TestExtractCodeBlocks_Single(t *testing.T) {
	input := "before ```\ncode here\n``` after"
	result := extractCodeBlocks(input)
	if len(result.codes) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(result.codes))
	}
	if !strings.Contains(result.codes[0], "code here") {
		t.Errorf("code block content = %q", result.codes[0])
	}
	if strings.Contains(result.text, "code here") {
		t.Error("text should have placeholder instead of code")
	}
}

func TestExtractCodeBlocks_Multiple(t *testing.T) {
	input := "```\nfirst\n``` middle ```\nsecond\n```"
	result := extractCodeBlocks(input)
	if len(result.codes) != 2 {
		t.Fatalf("expected 2 code blocks, got %d", len(result.codes))
	}
}

func TestExtractCodeBlocks_WithLanguage(t *testing.T) {
	input := "```go\nfunc main() {}\n```"
	result := extractCodeBlocks(input)
	if len(result.codes) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(result.codes))
	}
}

// --- extractInlineCodes ---

func TestExtractInlineCodes_None(t *testing.T) {
	result := extractInlineCodes("no inline code")
	if len(result.codes) != 0 {
		t.Errorf("expected 0, got %d", len(result.codes))
	}
}

func TestExtractInlineCodes_Single(t *testing.T) {
	result := extractInlineCodes("use `fmt.Println` for output")
	if len(result.codes) != 1 {
		t.Fatalf("expected 1, got %d", len(result.codes))
	}
	if result.codes[0] != "fmt.Println" {
		t.Errorf("code = %q, want %q", result.codes[0], "fmt.Println")
	}
}

func TestExtractInlineCodes_Multiple(t *testing.T) {
	result := extractInlineCodes("`a` and `b` and `c`")
	if len(result.codes) != 3 {
		t.Fatalf("expected 3, got %d", len(result.codes))
	}
}

// --- markdownToTelegramHTML ---

func TestMarkdownToTelegramHTML_Empty(t *testing.T) {
	got := markdownToTelegramHTML("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_Bold(t *testing.T) {
	got := markdownToTelegramHTML("**bold text**")
	if !strings.Contains(got, "<b>bold text</b>") {
		t.Errorf("expected bold HTML, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_Italic(t *testing.T) {
	got := markdownToTelegramHTML("_italic text_")
	if !strings.Contains(got, "<i>italic text</i>") {
		t.Errorf("expected italic HTML, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_Strikethrough(t *testing.T) {
	got := markdownToTelegramHTML("~~deleted~~")
	if !strings.Contains(got, "<s>deleted</s>") {
		t.Errorf("expected strikethrough HTML, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_Link(t *testing.T) {
	got := markdownToTelegramHTML("[click](https://example.com)")
	if !strings.Contains(got, `<a href="https://example.com">click</a>`) {
		t.Errorf("expected link HTML, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_Header(t *testing.T) {
	got := markdownToTelegramHTML("# Heading")
	// Headers are stripped to plain text
	if strings.Contains(got, "#") {
		t.Errorf("header marker should be removed, got %q", got)
	}
	if !strings.Contains(got, "Heading") {
		t.Errorf("heading text should be preserved, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_CodeBlock(t *testing.T) {
	got := markdownToTelegramHTML("```\ncode\n```")
	if !strings.Contains(got, "<pre><code>") {
		t.Errorf("expected code block HTML, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_InlineCode(t *testing.T) {
	got := markdownToTelegramHTML("use `fmt` package")
	if !strings.Contains(got, "<code>fmt</code>") {
		t.Errorf("expected inline code HTML, got %q", got)
	}
}

func TestMarkdownToTelegramHTML_HTMLEscaping(t *testing.T) {
	got := markdownToTelegramHTML("a < b & c > d")
	if !strings.Contains(got, "&lt;") || !strings.Contains(got, "&amp;") || !strings.Contains(got, "&gt;") {
		t.Errorf("HTML should be escaped, got %q", got)
	}
}

// --- parseChatID ---

func TestParseChatID_Valid(t *testing.T) {
	id, err := parseChatID("123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 123456 {
		t.Errorf("id = %d, want 123456", id)
	}
}

func TestParseChatID_Negative(t *testing.T) {
	id, err := parseChatID("-100123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != -100123 {
		t.Errorf("id = %d, want -100123", id)
	}
}

func TestParseChatID_Invalid(t *testing.T) {
	_, err := parseChatID("notanumber")
	if err == nil {
		t.Error("expected error for invalid chat ID")
	}
}

func TestParseChatID_Empty(t *testing.T) {
	_, err := parseChatID("")
	if err == nil {
		t.Error("expected error for empty chat ID")
	}
}
