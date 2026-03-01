package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewMemoryStore(t *testing.T) {
	tmpDir := t.TempDir()
	ms := NewMemoryStore(tmpDir)

	if ms == nil {
		t.Fatal("NewMemoryStore returned nil")
	}
	if ms.dataDir != tmpDir {
		t.Errorf("dataDir = %q, want %q", ms.dataDir, tmpDir)
	}

	// Verify memory directory was created
	memDir := filepath.Join(tmpDir, "memory")
	if _, err := os.Stat(memDir); os.IsNotExist(err) {
		t.Error("memory directory should be created")
	}
}

func TestReadLongTerm_NoFile(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())
	result := ms.ReadLongTerm()
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestWriteAndReadLongTerm(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())

	content := "This is long-term memory content."
	if err := ms.WriteLongTerm(content); err != nil {
		t.Fatalf("WriteLongTerm failed: %v", err)
	}

	got := ms.ReadLongTerm()
	if got != content {
		t.Errorf("ReadLongTerm = %q, want %q", got, content)
	}
}

func TestWriteLongTerm_Overwrite(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())

	_ = ms.WriteLongTerm("first")
	_ = ms.WriteLongTerm("second")

	got := ms.ReadLongTerm()
	if got != "second" {
		t.Errorf("ReadLongTerm = %q, want %q", got, "second")
	}
}

func TestReadToday_NoFile(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())
	result := ms.ReadToday()
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestAppendToday_NewFile(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())

	if err := ms.AppendToday("Today's note"); err != nil {
		t.Fatalf("AppendToday failed: %v", err)
	}

	got := ms.ReadToday()
	if !strings.Contains(got, "Today's note") {
		t.Errorf("expected content to contain 'Today's note', got %q", got)
	}
	// Should have date header
	today := time.Now().Format("2006-01-02")
	if !strings.Contains(got, today) {
		t.Errorf("expected date header %q in content", today)
	}
}

func TestAppendToday_AppendToExisting(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())

	_ = ms.AppendToday("first note")
	_ = ms.AppendToday("second note")

	got := ms.ReadToday()
	if !strings.Contains(got, "first note") {
		t.Error("should contain first note")
	}
	if !strings.Contains(got, "second note") {
		t.Error("should contain second note")
	}
}

func TestGetRecentDailyNotes_NoFiles(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())
	result := ms.GetRecentDailyNotes(3)
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestGetRecentDailyNotes_WithTodayNote(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())

	_ = ms.AppendToday("today's content")

	result := ms.GetRecentDailyNotes(3)
	if !strings.Contains(result, "today's content") {
		t.Errorf("expected today's content, got %q", result)
	}
}

func TestGetMemoryContext_Empty(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())
	result := ms.GetMemoryContext()
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestGetMemoryContext_LongTermOnly(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())
	_ = ms.WriteLongTerm("long-term info")

	result := ms.GetMemoryContext()
	if !strings.Contains(result, "# Memory") {
		t.Error("should contain Memory header")
	}
	if !strings.Contains(result, "Long-term Memory") {
		t.Error("should contain Long-term Memory section")
	}
	if !strings.Contains(result, "long-term info") {
		t.Error("should contain long-term content")
	}
}

func TestGetMemoryContext_BothSections(t *testing.T) {
	ms := NewMemoryStore(t.TempDir())
	_ = ms.WriteLongTerm("persistent info")
	_ = ms.AppendToday("daily note")

	result := ms.GetMemoryContext()
	if !strings.Contains(result, "Long-term Memory") {
		t.Error("should contain Long-term Memory section")
	}
	if !strings.Contains(result, "Recent Daily Notes") {
		t.Error("should contain Recent Daily Notes section")
	}
}
