// Package source provides the foundational data structures for tracking and
// managing source code text, locations, and diagnostics in the compiler.
package source

import (
	"bytes"
	"path/filepath"
	"slices"
	"unicode/utf8"
)

// File represents a single source code file tracked by the compiler.
// It caches the raw byte contents and pre-computes line offsets for fast
// O(log N) line and column lookups during tokenization and parsing.
type File struct {
	path        string
	version     int32
	text        []byte
	lineOffsets []uint32
}

// New initializes a new File instance, cloning the provided text buffer
// and pre-calculating all newline boundaries for fast positional lookups.
func New(path string, version int32, text []byte) *File {
	text = slices.Clone(text)

	lineCount := bytes.Count(text, []byte{'\n'})
	lineOffsets := make([]uint32, 0, lineCount+1)
	lineOffsets = append(lineOffsets, 0)

	var offset int
	for {
		idx := bytes.IndexByte(text[offset:], '\n')
		if idx < 0 {
			break
		}
		offset += idx + 1
		lineOffsets = append(lineOffsets, uint32(offset))
	}

	return &File{
		path:        path,
		version:     version,
		text:        text,
		lineOffsets: lineOffsets,
	}
}

// Metadata

// Path returns the absolute or relative file path.
func (f *File) Path() string { return f.path }

// Version returns the current version of the file, typically used by an LSP.
func (f *File) Version() int32 { return f.version }

// Base returns just the file name (e.g., "main.az").
func (f *File) Base() string { return filepath.Base(f.path) }

// Dir returns the directory containing the file.
func (f *File) Dir() string { return filepath.Dir(f.path) }

// Ext returns the file extension (e.g., ".az").
func (f *File) Ext() string { return filepath.Ext(f.path) }

// Contents

// String returns the complete file contents as a string.
func (f *File) String() string { return string(f.text) }

// Bytes returns a safe clone of the underlying byte slice.
func (f *File) Bytes() []byte { return slices.Clone(f.text) }

// Len returns the total length of the file in bytes.
func (f *File) Len() uint32 { return uint32(len(f.text)) }

// Empty reports whether the file has a length of zero.
func (f *File) Empty() bool { return len(f.text) == 0 }

// Offsets

// Valid reports whether the given offset is within the file bounds.
func (f *File) Valid(offset uint32) bool { return offset <= f.Len() }

// EOF reports whether the offset is at or beyond the end of the file.
func (f *File) EOF(offset uint32) bool { return offset >= f.Len() }

// Clamp constrains the given offset to be at most the file's maximum length.
func (f *File) Clamp(offset uint32) uint32 {
	if offset > f.Len() {
		return f.Len()
	}
	return offset
}

// Lines

// LineCount returns the total number of lines in the file.
func (f *File) LineCount() uint32 { return uint32(len(f.lineOffsets)) }

// Line retrieves the physical line information for a 1-based line number.
// It excludes the terminating \r or \n characters.
func (f *File) Line(number uint32) Line {
	if number == 0 || number > f.LineCount() {
		return Line{}
	}

	start := f.lineOffsets[number-1]
	end := f.Len()
	if number < f.LineCount() {
		end = f.lineOffsets[number]
	}

	for end > start {
		switch f.text[end-1] {
		case '\r', '\n':
			end--
		default:
			goto done
		}
	}

done:
	return Line{
		Number: number,
		Start:  start,
		End:    end,
	}
}

// LineText returns the raw bytes for the specified 1-based line number,
// excluding the trailing newline characters.
func (f *File) LineText(number uint32) []byte {
	line := f.Line(number)
	if line.Empty() {
		return nil
	}
	return f.text[line.Start:line.End]
}

// searchLineIndex performs a fast O(log N) binary search to find the 0-based
// index of the line containing the given byte offset.
func (f *File) searchLineIndex(offset uint32) int {
	offset = f.Clamp(offset)

	low, high := 0, len(f.lineOffsets)
	for low < high {
		mid := low + (high-low)>>1
		if f.lineOffsets[mid] <= offset {
			low = mid + 1
		} else {
			high = mid
		}
	}

	if low == 0 {
		return 0
	}
	return low - 1
}

// LineAt returns the physical line information enclosing the given byte offset.
func (f *File) LineAt(offset uint32) Line {
	return f.Line(uint32(f.searchLineIndex(offset) + 1))
}

// Positions

// Position translates a flat byte offset into a structural Position
// containing the 1-based line, byte column, and LSP-compliant column.
func (f *File) Position(offset uint32) Position {
	offset = f.Clamp(offset)
	lineIdx := f.searchLineIndex(offset)

	lineStart := f.lineOffsets[lineIdx]
	byteColumn := offset - lineStart

	var lspColumn uint32
	prefixText := f.text[lineStart:offset]
	for len(prefixText) > 0 {
		r, size := utf8.DecodeRune(prefixText)
		if r >= 0x10000 {
			lspColumn += 2
		} else {
			lspColumn += 1
		}
		prefixText = prefixText[size:]
	}

	return Position{
		Offset:     offset,
		Line:       uint32(lineIdx + 1),
		ByteColumn: byteColumn + 1,
		LSPColumn:  lspColumn,
	}
}

// Offset converts a 1-based line and 1-based column into a flat byte offset.
func (f *File) Offset(line, column uint32) uint32 {
	l := f.Line(line)
	if l.Number == 0 {
		return f.Len()
	}

	offset := l.Start + column - 1
	if offset > l.End {
		return l.End
	}

	return offset
}

// Spans

// BytesOf returns a slice of the file's text bounded by the given Span.
func (f *File) BytesOf(span Span) []byte {
	start := f.Clamp(span.Start)
	end := f.Clamp(span.End)
	if end < start {
		end = start
	}
	return f.text[start:end]
}

// Text returns the string representation of the text bounded by the given Span.
func (f *File) Text(span Span) string {
	return string(f.BytesOf(span))
}

// LineSpan returns a Span covering the entirety of the specified 1-based line.
func (f *File) LineSpan(number uint32) Span {
	line := f.Line(number)
	return NewSpan(line.Start, line.End)
}

// UTF-8

// Rune decodes and returns the rune at the given offset, along with its size in bytes.
func (f *File) Rune(offset uint32) (rune, uint32) {
	offset = f.Clamp(offset)
	if offset >= f.Len() {
		return utf8.RuneError, 0
	}
	r, size := utf8.DecodeRune(f.text[offset:])
	return r, uint32(size)
}

// RuneBefore decodes and returns the rune immediately preceding the given offset.
func (f *File) RuneBefore(offset uint32) (rune, uint32) {
	offset = f.Clamp(offset)
	if offset == 0 {
		return utf8.RuneError, 0
	}
	r, size := utf8.DecodeLastRune(f.text[:offset])
	return r, uint32(size)
}
