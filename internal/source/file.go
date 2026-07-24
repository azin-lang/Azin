// Package source handles the reading, slicing, and positional indexing of source files.
package source

import (
	"fmt"
	"path/filepath"
	"sort"
	"unicode/utf8"

	"github.com/azin-lang/Azin/internal/token"
)

// File holds the contents and line structure of a parsed source file.
type File struct {
	name  string
	text  []byte
	lines []uint32
}

// New returns a new File, parsing line offsets from the provided text.
func New(name string, text []byte) *File {
	lines := []uint32{0}

	for i, ch := range text {
		if ch == '\n' {
			lines = append(lines, uint32(i+1))
		}
	}

	return &File{
		name:  name,
		text:  text,
		lines: lines,
	}
}

// Len returns the size of the file in bytes.
func (f *File) Len() uint32 {
	return uint32(len(f.text)) //nolint:gosec
}

// Empty reports whether the file contains no text.
func (f *File) Empty() bool {
	return f.Len() == 0
}

// EOF reports whether the offset is at or beyond the end of the file.
func (f *File) EOF(offset uint32) bool {
	return offset >= f.Len()
}

// Byte returns the character at the given offset.
func (f *File) Byte(offset uint32) byte {
	return f.text[offset]
}

// Rune returns the rune at the given offset, along with the byte count
func (f *File) Rune(offset uint32) (r rune, size uint32) {
	var sz int
	r, sz = utf8.DecodeRune(f.text[offset:])
	return r, uint32(sz) //nolint:gosec
}

// Slice returns the byte slice between start and end offsets.
func (f *File) Slice(start, end uint32) []byte {
	return f.text[start:end]
}

// Text returns the exact byte slice for a given token.
func (f *File) Text(tok token.Token) []byte {
	return f.Slice(tok.Position.Offset, tok.Position.Offset+tok.Length)
}

// LineColumn returns the 1-based line and column numbers for the offset.
func (f *File) LineColumn(offset uint32) (line, column uint32) {
	var i = max(0, sort.Search(len(f.lines), func(i int) bool {
		return f.lines[i] > offset
	})-1)

	line = uint32(i + 1) //nolint:gosec
	column = offset - f.lines[i] + 1
	return
}

// LineCount returns the total number of lines in the file.
func (f *File) LineCount() uint32 {
	return uint32(len(f.lines)) //nolint:gosec
}

// LineStart returns the byte offset for the beginning of a 1-based line.
func (f *File) LineStart(line uint32) (uint32, bool) {
	if line == 0 || line > uint32(len(f.lines)) { //nolint:gosec
		return 0, false
	}

	return f.lines[line-1], true
}

// Line returns the byte slice for a 1-based line, excluding trailing newlines.
func (f *File) Line(line uint32) []byte {
	start, ok := f.LineStart(line)
	if !ok {
		return nil
	}

	var end uint32
	if line == f.LineCount() {
		end = f.Len()
	} else {
		end = f.lines[line] // start of next line
	}

	for end > start {
		switch f.text[end-1] {
		case '\n', '\r':
			end--
		default:
			return f.text[start:end]
		}
	}

	return f.text[start:end]
}

// LineText returns the full text of the line containing the given offset.
func (f *File) LineText(offset uint32) []byte {
	line, _ := f.LineColumn(offset)
	return f.Line(line)
}

// Name returns the file path.
func (f *File) Name() string {
	return f.name
}

// Base returns the base file name.
func (f *File) Base() string {
	return filepath.Base(f.name)
}

// Dir returns the directory portion of the file path.
func (f *File) Dir() string {
	return filepath.Dir(f.name)
}

// Ext returns the file extension.
func (f *File) Ext() string {
	return filepath.Ext(f.name)
}

// FormatToken returns a string representation of a token, including its location and text.
func (f *File) FormatToken(tok token.Token) string {
	line, column := f.LineColumn(tok.Position.Offset)

	s := fmt.Sprintf(
		"%-18s %4d:%4d [%d:%d]",
		tok.Kind,
		line,
		column,
		tok.Position.Offset,
		tok.Length,
	)

	if tok.Kind.HasText() {
		s += fmt.Sprintf(" %q", f.Text(tok))
	}

	return s
}
