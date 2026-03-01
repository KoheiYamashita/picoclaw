package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normal filename", "photo.jpg", "photo.jpg"},
		{"path traversal", "../../../etc/passwd", "passwd"},
		{"forward slash", "path/to/file.txt", "file.txt"},
		{"backslash", "path\\to\\file.txt", "path_to_file.txt"},
		{"double dot", "file..name.txt", "filename.txt"},
		{"empty becomes dot", "", "."},
		{"unicode filename", "日本語ファイル.png", "日本語ファイル.png"},
		{"spaces in name", "my file.txt", "my file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFilename(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEncodeFileToDataURL_JPEG(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.jpg")
	if err := os.WriteFile(path, []byte("fake jpeg data"), 0644); err != nil {
		t.Fatal(err)
	}

	result := EncodeFileToDataURL(path)
	if result == "" {
		t.Error("expected non-empty data URL for .jpg")
	}
	if len(result) < 20 {
		t.Error("data URL seems too short")
	}
	if result[:15] != "data:image/jpeg" {
		t.Errorf("expected jpeg MIME type, got prefix %q", result[:15])
	}
}

func TestEncodeFileToDataURL_PNG(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.png")
	if err := os.WriteFile(path, []byte("fake png data"), 0644); err != nil {
		t.Fatal(err)
	}

	result := EncodeFileToDataURL(path)
	if result == "" {
		t.Error("expected non-empty data URL for .png")
	}
	if len(result) < 14 {
		t.Error("data URL seems too short")
	}
	if result[:14] != "data:image/png" {
		t.Errorf("expected png MIME type, got prefix %q", result[:14])
	}
}

func TestEncodeFileToDataURL_Unsupported(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(path, []byte("text content"), 0644); err != nil {
		t.Fatal(err)
	}

	result := EncodeFileToDataURL(path)
	if result != "" {
		t.Errorf("expected empty string for unsupported extension, got %q", result)
	}
}

func TestEncodeFileToDataURL_NonExistent(t *testing.T) {
	result := EncodeFileToDataURL("/nonexistent/file.jpg")
	if result != "" {
		t.Errorf("expected empty string for non-existent file, got %q", result)
	}
}

func TestEncodeFileToDataURL_TooLarge(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "large.jpg")

	// Create a file larger than 50MB
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	// Write just the header to fake the size - use Truncate to set size
	if err := f.Truncate(51 * 1024 * 1024); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	_ = f.Close()

	result := EncodeFileToDataURL(path)
	if result != "" {
		t.Error("expected empty string for oversized file")
	}
}

func TestEncodeFileToDataURL_AllSupportedExtensions(t *testing.T) {
	tmpDir := t.TempDir()
	extensions := []struct {
		ext  string
		mime string
	}{
		{".jpg", "image/jpeg"},
		{".jpeg", "image/jpeg"},
		{".png", "image/png"},
		{".webp", "image/webp"},
		{".gif", "image/gif"},
	}

	for _, ext := range extensions {
		t.Run(ext.ext, func(t *testing.T) {
			path := filepath.Join(tmpDir, "test"+ext.ext)
			if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
				t.Fatal(err)
			}
			result := EncodeFileToDataURL(path)
			if result == "" {
				t.Errorf("expected non-empty data URL for %s", ext.ext)
			}
		})
	}
}

// --- DownloadFile ---

func TestDownloadFile_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("file content"))
	}))
	defer ts.Close()

	result := DownloadFile(ts.URL, "test.txt", DownloadOptions{})
	if result == "" {
		t.Fatal("expected non-empty path")
	}
	defer func() { _ = os.Remove(result) }()

	data, err := os.ReadFile(result)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if string(data) != "file content" {
		t.Errorf("content = %q, want %q", string(data), "file content")
	}
}

func TestDownloadFile_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	result := DownloadFile(ts.URL, "test.txt", DownloadOptions{})
	if result != "" {
		_ = os.Remove(result)
		t.Error("expected empty string for 404 response")
	}
}

func TestDownloadFile_ExtraHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token123" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	result := DownloadFile(ts.URL, "test.txt", DownloadOptions{
		ExtraHeaders: map[string]string{
			"Authorization": "Bearer token123",
		},
	})
	if result == "" {
		t.Fatal("expected non-empty path with correct auth header")
	}
	defer func() { _ = os.Remove(result) }()
}

func TestDownloadFile_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	result := DownloadFile(ts.URL, "test.txt", DownloadOptions{
		Timeout: 100 * time.Millisecond,
	})
	if result != "" {
		_ = os.Remove(result)
		t.Error("expected empty string for timeout")
	}
}

func TestDownloadFileSimple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("simple content"))
	}))
	defer ts.Close()

	result := DownloadFileSimple(ts.URL, "simple.txt")
	if result == "" {
		t.Fatal("expected non-empty path")
	}
	defer func() { _ = os.Remove(result) }()

	data, err := os.ReadFile(result)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "simple content" {
		t.Errorf("content = %q", string(data))
	}
}
