package source_test

import (
	"path/filepath"
	"testing"

	fs "github.com/azin-lang/Azin/pkg/source"
	"github.com/azin-lang/Azin/pkg/token"
)

func TestNewFile(t *testing.T) {
	f := fs.New("test.az", []byte("hello\nworld\n"))
	if f.Name() != "test.az" {
		t.Errorf("Name() = %q, want %q", f.Name(), "test.az")
	}
	if f.Len() != 12 {
		t.Errorf("Len() = %d, want 12", f.Len())
	}
	if f.Empty() {
		t.Error("Empty() = true, want false")
	}
}

func TestEmptyFile(t *testing.T) {
	f := fs.New("empty.az", nil)
	if !f.Empty() {
		t.Error("Empty() = false, want true")
	}
	if f.Len() != 0 {
		t.Errorf("Len() = %d, want 0", f.Len())
	}
	if !f.EOF(0) {
		t.Error("EOF(0) = false, want true")
	}
}

func TestLineColumn(t *testing.T) {
	f := fs.New("test.az", []byte("line1\nline2\nline3"))
	tests := []struct {
		offset   uint32
		wantLine uint32
		wantCol  uint32
	}{
		{0, 1, 1},
		{5, 1, 6},
		{6, 2, 1},
		{11, 2, 6},
		{12, 3, 1},
		{16, 3, 5},
	}

	for _, tt := range tests {
		line, col := f.LineColumn(tt.offset)
		if line != tt.wantLine || col != tt.wantCol {
			t.Errorf("LineColumn(%d) = (%d,%d), want (%d,%d)",
				tt.offset, line, col, tt.wantLine, tt.wantCol)
		}
	}
}

func TestLine(t *testing.T) {
	f := fs.New("test.az", []byte("abc\ndef\nghi"))
	if got := string(f.Line(1)); got != "abc" {
		t.Errorf("Line(1) = %q, want %q", got, "abc")
	}
	if got := string(f.Line(2)); got != "def" {
		t.Errorf("Line(2) = %q, want %q", got, "def")
	}
	if got := string(f.Line(3)); got != "ghi" {
		t.Errorf("Line(3) = %q, want %q", got, "ghi")
	}
}

func TestLineOutOfRange(t *testing.T) {
	f := fs.New("test.az", []byte("abc"))
	if got := f.Line(0); got != nil {
		t.Errorf("Line(0) = %v, want nil", got)
	}
	if got := f.Line(99); got != nil {
		t.Errorf("Line(99) = %v, want nil", got)
	}
}

func TestRune(t *testing.T) {
	f := fs.New("test.az", []byte("a∂c"))
	r, size := f.Rune(0)
	if r != 'a' || size != 1 {
		t.Errorf("Rune(0) = (%c,%d), want (%c,%d)", r, size, 'a', 1)
	}
	r, size = f.Rune(1)
	if r != '∂' || size != 3 {
		t.Errorf("Rune(1) = (%c,%d), want (%c,%d)", r, size, '∂', 3)
	}
}

func TestSlice(t *testing.T) {
	f := fs.New("test.az", []byte("hello world"))
	if got := string(f.Slice(0, 5)); got != "hello" {
		t.Errorf("Slice(0,5) = %q, want %q", got, "hello")
	}
	if got := string(f.Slice(6, 11)); got != "world" {
		t.Errorf("Slice(6,11) = %q, want %q", got, "world")
	}
}

func TestText(t *testing.T) {
	f := fs.New("test.az", []byte("var x: int"))
	tok := token.Token{Kind: token.Identifier, Position: token.Position{Offset: 4}, Length: 1}
	if got := string(f.Text(tok)); got != "x" {
		t.Errorf("Text(identifier) = %q, want %q", got, "x")
	}
}

func TestBaseAndExt(t *testing.T) {
	f := fs.New("/path/to/file.az", nil)
	if f.Base() != "file.az" {
		t.Errorf("Base() = %q, want %q", f.Base(), "file.az")
	}
	if f.Ext() != ".az" {
		t.Errorf("Ext() = %q, want %q", f.Ext(), ".az")
	}
	wantDir := filepath.FromSlash("/path/to")
	if f.Dir() != wantDir {
		t.Errorf("Dir() = %q, want %q", f.Dir(), wantDir)
	}
}

func TestLineCount(t *testing.T) {
	tests := []struct {
		input string
		want  uint32
	}{
		{"", 1},
		{"hello", 1},
		{"hello\n", 2},
		{"hello\nworld", 2},
		{"hello\nworld\n", 3},
	}

	for _, tt := range tests {
		f := fs.New("test.az", []byte(tt.input))
		if got := f.LineCount(); got != tt.want {
			t.Errorf("LineCount(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestEOF(t *testing.T) {
	f := fs.New("test.az", []byte("abc"))
	tests := []struct {
		offset uint32
		want   bool
	}{
		{0, false},
		{2, false},
		{3, true},
		{100, true},
	}

	for _, tt := range tests {
		if got := f.EOF(tt.offset); got != tt.want {
			t.Errorf("EOF(%d) = %v, want %v", tt.offset, got, tt.want)
		}
	}
}

func TestFormatToken(t *testing.T) {
	f := fs.New("test.az", []byte("fn main()"))
	tok := token.Token{
		Kind:     token.Identifier,
		Position: token.Position{Offset: 3},
		Length:   4,
	}
	got := f.FormatToken(tok)
	if got == "" {
		t.Error("FormatToken returned empty string")
	}
}
