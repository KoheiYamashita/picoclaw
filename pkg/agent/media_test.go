package agent

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/KarakuriAgent/clawdroid/pkg/providers"
)

// --- mimeToExt ---

func TestMimeToExt(t *testing.T) {
	tests := []struct {
		mime string
		want string
	}{
		{"image/jpeg", ".jpg"},
		{"image/png", ".png"},
		{"image/gif", ".gif"},
		{"image/webp", ".webp"},
		{"image/bmp", ".bmp"},
		{"application/octet-stream", ".bin"},
		{"unknown/type", ".bin"},
		{"IMAGE/JPEG", ".jpg"}, // case-insensitive
	}

	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			got := mimeToExt(tt.mime)
			if got != tt.want {
				t.Errorf("mimeToExt(%q) = %q, want %q", tt.mime, got, tt.want)
			}
		})
	}
}

// --- parseDataURL ---

func TestParseDataURL_Valid(t *testing.T) {
	data := []byte("hello world")
	encoded := base64.StdEncoding.EncodeToString(data)
	dataURL := "data:image/png;base64," + encoded

	ext, decoded, err := parseDataURL(dataURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ext != ".png" {
		t.Errorf("ext = %q, want %q", ext, ".png")
	}
	if string(decoded) != "hello world" {
		t.Errorf("decoded = %q, want %q", string(decoded), "hello world")
	}
}

func TestParseDataURL_NotDataURL(t *testing.T) {
	_, _, err := parseDataURL("https://example.com/image.png")
	if err == nil {
		t.Error("expected error for non-data URL")
	}
}

func TestParseDataURL_NoComma(t *testing.T) {
	_, _, err := parseDataURL("data:image/png;base64")
	if err == nil {
		t.Error("expected error for data URL without comma")
	}
}

func TestParseDataURL_InvalidBase64(t *testing.T) {
	_, _, err := parseDataURL("data:image/png;base64,!!!invalid!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

// --- PersistMedia ---

func TestPersistMedia_Empty(t *testing.T) {
	result := PersistMedia(nil, t.TempDir())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestPersistMedia_EmptyMediaDir(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("test"))
	result := PersistMedia([]string{"data:image/png;base64," + encoded}, "")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestPersistMedia_ValidDataURL(t *testing.T) {
	tmpDir := t.TempDir()
	data := []byte("fake image data")
	encoded := base64.StdEncoding.EncodeToString(data)
	dataURL := "data:image/png;base64," + encoded

	result := PersistMedia([]string{dataURL}, tmpDir)
	if len(result) != 1 {
		t.Fatalf("expected 1 path, got %d", len(result))
	}

	// Verify file exists and has correct content
	content, err := os.ReadFile(result[0])
	if err != nil {
		t.Fatalf("failed to read persisted file: %v", err)
	}
	if string(content) != "fake image data" {
		t.Errorf("file content = %q, want %q", string(content), "fake image data")
	}
}

func TestPersistMedia_SkipNonDataURLs(t *testing.T) {
	tmpDir := t.TempDir()
	result := PersistMedia([]string{"/local/path/file.png", "https://example.com/image.png"}, tmpDir)
	if len(result) != 0 {
		t.Errorf("expected 0 paths for non-data URLs, got %d", len(result))
	}
}

// --- CleanupMediaFiles ---

func TestCleanupMediaFiles_EmptyMessages(t *testing.T) {
	// Should not panic
	CleanupMediaFiles(nil)
	CleanupMediaFiles([]providers.Message{})
}

func TestCleanupMediaFiles_WithImagePaths(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "image.jpg")
	if err := os.WriteFile(tmpFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	messages := []providers.Message{
		{Content: "Here is [Image: " + tmpFile + "] for you"},
	}

	CleanupMediaFiles(messages)

	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("image file should have been deleted")
	}
}

func TestCleanupMediaFiles_NonExistentFile(t *testing.T) {
	// Should not panic on non-existent files
	messages := []providers.Message{
		{Content: "[Image: /nonexistent/file.jpg]"},
	}
	CleanupMediaFiles(messages)
}
