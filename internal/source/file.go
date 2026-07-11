package source

import (
	"fmt"           // fmt package is used for formatting strings
	"path/filepath" // filepath package is used for manipulating file paths
	"sort"          // sort package is used for sorting slices

	"github.com/azin-lang/Azin/internal/token" // token package is used for tokenizing source code
)

// File represents a source file.
type File struct {
	name  string
	text  []byte
	lines []uint32
}

// New returns a new File.
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

// Len returns the length of the file.
func (f *File) Len() uint32 {
	return uint32(len(f.text))
}

// Empty returns whether the file is empty.
func (f *File) Empty() bool {
	return f.Len() == 0
}

// EOF returns whether the offset is past the end of the file.
func (f *File) EOF(offset uint32) bool {
	return offset >= f.Len()
}

// Byte returns the byte at the given offset.
func (f *File) Byte(offset uint32) byte {
	return f.text[offset]
}

// Slice returns a slice of the file's text.
func (f *File) Slice(start, end uint32) []byte {
	return f.text[start:end]
}

// Text returns the text of a token.
func (f *File) Text(tok token.Token) []byte {
	return f.Slice(tok.Position.Offset, tok.Position.Offset+tok.Length)
}

// LineColumn returns the line and column of an offset.
func (f *File) LineColumn(offset uint32) (line, column uint32) {
	var i = max(0, sort.Search(len(f.lines), func(i int) bool {
		return f.lines[i] > offset
	})-1)

	line = uint32(i + 1)
	column = offset - f.lines[i] + 1
	return
}

// LineCount returns the number of lines in the file.
func (f *File) LineCount() uint32 {
	return uint32(len(f.lines))
}

// LineStart returns the start offset of a line.
func (f *File) LineStart(line uint32) (uint32, bool) {
	if line == 0 || line > uint32(len(f.lines)) {
		return 0, false
	}

	return f.lines[line-1], true
}

// Line returns the text of a line.
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

	// Trim trailing line endings.
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

// LineText returns the text of a line.
func (f *File) LineText(offset uint32) []byte {
	line, _ := f.LineColumn(offset)
	return f.Line(line)
}

// Name returns the name of the file.
func (f *File) Name() string {
	return f.name
}

// Base returns the base name of the file.
func (f *File) Base() string {
	return filepath.Base(f.name)
}

// Dir returns the directory of the file.
func (f *File) Dir() string {
	return filepath.Dir(f.name)
}

// Ext returns the extension of the file.
func (f *File) Ext() string {
	return filepath.Ext(f.name)
}

// FormatToken formats a token as a string.
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
