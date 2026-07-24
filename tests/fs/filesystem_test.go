package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/azin-lang/Azin/internal/fs"
)

func TestReadSourceFileValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.az")
	if err := os.WriteFile(path, []byte("fn main() do end"), 0o600); err != nil {
		t.Fatal(err)
	}

	data, err := fs.ReadSourceFile(path, false)
	if err != nil {
		t.Fatalf("ReadSourceFile(%q) = %v", path, err)
	}
	if string(data) != "fn main() do end" {
		t.Errorf("got %q, want %q", string(data), "fn main() do end")
	}
}

func TestReadSourceFileInvalidExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := fs.ReadSourceFile(path, false)
	if err == nil {
		t.Error("expected error for .txt file, got nil")
	}
}

func TestReadSourceFileIgnoreExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}

	data, err := fs.ReadSourceFile(path, true)
	if err != nil {
		t.Fatalf("ReadSourceFile with ignoreExtension: %v", err)
	}
	if string(data) != "content" {
		t.Errorf("got %q, want %q", string(data), "content")
	}
}

func TestReadSourceFileNotFound(t *testing.T) {
	_, err := fs.ReadSourceFile("/nonexistent/path.az", false)
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestReadSourceFileValidExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.az")
	if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}

	data, err := fs.ReadSourceFile(path, false)
	if err != nil {
		t.Errorf("ReadSourceFile(%q): %v", path, err)
	}
	if string(data) != "content" {
		t.Errorf("ReadSourceFile(%q) = %q, want %q", path, string(data), "content")
	}
}

func TestReadSourceFileEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.az")
	if err := os.WriteFile(path, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}

	data, err := fs.ReadSourceFile(path, false)
	if err != nil {
		t.Fatalf("ReadSourceFile(%q): %v", path, err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty data, got %d bytes", len(data))
	}
}
