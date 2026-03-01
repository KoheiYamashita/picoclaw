package channels

import (
	"strings"
	"testing"
)

// --- splitMessage ---

func TestSplitMessage_ShortMessage(t *testing.T) {
	msg := "hello world"
	chunks := splitMessage(msg, 100)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0] != msg {
		t.Errorf("chunk = %q, want %q", chunks[0], msg)
	}
}

func TestSplitMessage_EmptyString(t *testing.T) {
	chunks := splitMessage("", 100)
	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for empty string, got %d", len(chunks))
	}
}

func TestSplitMessage_ExactLimit(t *testing.T) {
	msg := strings.Repeat("a", 100)
	chunks := splitMessage(msg, 100)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
}

func TestSplitMessage_SplitAtNewline(t *testing.T) {
	// Build a message that is > limit with a newline near the end
	part1 := strings.Repeat("a", 80)
	part2 := strings.Repeat("b", 80)
	msg := part1 + "\n" + part2
	chunks := splitMessage(msg, 100)
	if len(chunks) < 2 {
		t.Fatalf("expected at least 2 chunks, got %d", len(chunks))
	}
	if chunks[0] != part1 {
		t.Errorf("first chunk = %q, want %q", chunks[0], part1)
	}
	// Verify no content is lost
	joined := strings.Join(chunks, "\n")
	if joined != msg {
		t.Errorf("content lost: joined length=%d, original length=%d", len(joined), len(msg))
	}
}

func TestSplitMessage_SplitAtSpace(t *testing.T) {
	// No newlines, but has spaces
	part1 := strings.Repeat("a", 80)
	part2 := strings.Repeat("b", 80)
	msg := part1 + " " + part2
	chunks := splitMessage(msg, 100)
	if len(chunks) < 2 {
		t.Fatalf("expected at least 2 chunks, got %d", len(chunks))
	}
	if chunks[0] != part1 {
		t.Errorf("first chunk = %q, want %q", chunks[0], part1)
	}
	if chunks[1] != part2 {
		t.Errorf("second chunk = %q, want %q", chunks[1], part2)
	}
}

func TestSplitMessage_CodeBlockPreserved(t *testing.T) {
	// Code block that starts before limit and closes after
	code := "```\n" + strings.Repeat("x", 120) + "\n```"
	msg := "intro\n" + code
	chunks := splitMessage(msg, 100)

	// The code block should be kept intact (extended)
	joined := strings.Join(chunks, "")
	if !strings.Contains(joined, "```") {
		t.Error("code block markers should be preserved")
	}
}

// --- findLastUnclosedCodeBlock ---

func TestFindLastUnclosedCodeBlock_NoCodeBlock(t *testing.T) {
	idx := findLastUnclosedCodeBlock("no code here")
	if idx != -1 {
		t.Errorf("expected -1, got %d", idx)
	}
}

func TestFindLastUnclosedCodeBlock_Balanced(t *testing.T) {
	text := "before ```code``` after"
	idx := findLastUnclosedCodeBlock(text)
	if idx != -1 {
		t.Errorf("expected -1 for balanced code block, got %d", idx)
	}
}

func TestFindLastUnclosedCodeBlock_Unbalanced(t *testing.T) {
	text := "before ```code here"
	idx := findLastUnclosedCodeBlock(text)
	if idx != 7 {
		t.Errorf("expected index 7, got %d", idx)
	}
}

func TestFindLastUnclosedCodeBlock_MultiplePairs(t *testing.T) {
	text := "```a``` ```b```"
	idx := findLastUnclosedCodeBlock(text)
	if idx != -1 {
		t.Errorf("expected -1 for multiple balanced pairs, got %d", idx)
	}
}

func TestFindLastUnclosedCodeBlock_OddCount(t *testing.T) {
	text := "```a``` ```b"
	idx := findLastUnclosedCodeBlock(text)
	// Implementation tracks lastOpenIdx as the first ``` when count transitions from 0.
	// With 3 markers total (odd), it returns the position of the first opening marker.
	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}
}

// --- findNextClosingCodeBlock ---

func TestFindNextClosingCodeBlock_NotFound(t *testing.T) {
	text := "no code block here"
	idx := findNextClosingCodeBlock(text, 0)
	if idx != -1 {
		t.Errorf("expected -1, got %d", idx)
	}
}

func TestFindNextClosingCodeBlock_Found(t *testing.T) {
	text := "some text ```after"
	idx := findNextClosingCodeBlock(text, 0)
	if idx != 13 {
		t.Errorf("expected index 13 (position after ```), got %d", idx)
	}
}

func TestFindNextClosingCodeBlock_AfterStartIdx(t *testing.T) {
	text := "```first``` and ```second```"
	// Start searching after the first closing
	idx := findNextClosingCodeBlock(text, 12)
	if idx <= 12 {
		t.Errorf("expected index > 12, got %d", idx)
	}
}

// --- findLastNewline ---

func TestFindLastNewline_Found(t *testing.T) {
	text := "line1\nline2\nline3"
	idx := findLastNewline(text, 200)
	if idx != 11 {
		t.Errorf("expected 11, got %d", idx)
	}
}

func TestFindLastNewline_NotFound(t *testing.T) {
	text := "no newlines here"
	idx := findLastNewline(text, 200)
	if idx != -1 {
		t.Errorf("expected -1, got %d", idx)
	}
}

func TestFindLastNewline_WindowConstraint(t *testing.T) {
	// Newline is outside the search window
	text := "line1\n" + strings.Repeat("a", 100)
	idx := findLastNewline(text, 10) // Only search last 10 chars
	if idx != -1 {
		t.Errorf("expected -1 (newline outside window), got %d", idx)
	}
}

// --- findLastSpace ---

func TestFindLastSpace_Found(t *testing.T) {
	text := "word1 word2 word3"
	idx := findLastSpace(text, 200)
	if idx != 11 {
		t.Errorf("expected 11, got %d", idx)
	}
}

func TestFindLastSpace_Tab(t *testing.T) {
	text := "word1\tword2"
	idx := findLastSpace(text, 200)
	if idx != 5 {
		t.Errorf("expected 5, got %d", idx)
	}
}

func TestFindLastSpace_WindowConstraint(t *testing.T) {
	text := "word " + strings.Repeat("x", 100)
	idx := findLastSpace(text, 10)
	if idx != -1 {
		t.Errorf("expected -1 (space outside window), got %d", idx)
	}
}

func TestFindLastSpace_NotFound(t *testing.T) {
	text := "nospaces"
	idx := findLastSpace(text, 200)
	if idx != -1 {
		t.Errorf("expected -1, got %d", idx)
	}
}

// --- appendContent ---

func TestAppendContent_EmptyContent(t *testing.T) {
	result := appendContent("", "suffix")
	if result != "suffix" {
		t.Errorf("expected %q, got %q", "suffix", result)
	}
}

func TestAppendContent_EmptySuffix(t *testing.T) {
	result := appendContent("content", "")
	if result != "content\n" {
		t.Errorf("expected %q, got %q", "content\n", result)
	}
}

func TestAppendContent_BothPresent(t *testing.T) {
	result := appendContent("hello", "world")
	expected := "hello\nworld"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestAppendContent_BothEmpty(t *testing.T) {
	result := appendContent("", "")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
