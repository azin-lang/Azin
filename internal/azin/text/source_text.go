// Package source provides the foundational data structures for tracking and
// managing source code text, locations, and diagnostics in the compiler.
package text

import (
	"bytes"
	"path/filepath"
	"slices"
	"unicode/utf8"
)

// SourceText represents a single source code file tracked by the compiler.
// It caches the raw byte contents and pre-computes line offsets for fast
// O(log N) line and column lookups during tokenization and parsing.
type SourceText struct {
	path        string
	version     int32
	text        []byte
	lineOffsets []int
}

// NewSourceText initializes a new SourceText instance, cloning the provided text buffer
// and pre-calculating all newline boundaries for fast positional lookups.
func NewSourceText(path string, version int32, text []byte) *SourceText {
	text = slices.Clone(text)

	lineCount := bytes.Count(text, []byte{'\n'})
	lineOffsets := make([]int, 0, lineCount+1)
	lineOffsets = append(lineOffsets, 0)

	offset := 0
	for {
		idx := bytes.IndexByte(text[offset:], '\n')
		if idx < 0 {
			break
		}

		offset += idx + 1
		lineOffsets = append(lineOffsets, offset)
	}

	return &SourceText{
		path:        path,
		version:     version,
		text:        text,
		lineOffsets: lineOffsets,
	}
}

// Metadata

// Path returns the absolute or relative file path.
func (f *SourceText) Path() string { return f.path }

// Version returns the current version of the file, typically used by an LSP.
func (f *SourceText) Version() int32 { return f.version }

// Base returns just the file name (e.g., "main.az").
func (f *SourceText) Base() string { return filepath.Base(f.path) }

// Dir returns the directory containing the file.
func (f *SourceText) Dir() string { return filepath.Dir(f.path) }

// Ext returns the file extension (e.g., ".az").
func (f *SourceText) Ext() string { return filepath.Ext(f.path) }

// Contents

// String returns the complete file contents as a string.
func (f *SourceText) String() string { return string(f.text) }

// Bytes returns a safe clone of the underlying byte slice.
func (f *SourceText) Bytes() []byte { return slices.Clone(f.text) }

// Len returns the total length of the file in bytes.
func (f *SourceText) Len() int { return len(f.text) }

// Empty reports whether the file has a length of zero.
func (f *SourceText) Empty() bool { return len(f.text) == 0 }

// Offsets

// Valid reports whether the given offset is within the file bounds.
func (f *SourceText) Valid(offset int) bool {
	return offset >= 0 && offset <= f.Len()
}

// EOF reports whether the offset is at or beyond the end of the file.
func (f *SourceText) EOF(offset int) bool {
	return offset >= f.Len()
}

// Clamp constrains the given offset to the valid file bounds.
func (f *SourceText) Clamp(offset int) int {
	return max(0, min(offset, f.Len()))
}

// Lines

// LineCount returns the total number of lines in the file.
func (f *SourceText) LineCount() int { return len(f.lineOffsets) }

// Line retrieves the physical line information for a 1-based line number.
// It excludes the terminating \r or \n characters.
func (f *SourceText) Line(number int) Line {
	if number <= 0 || number > f.LineCount() {
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
		Span:   SpanFromBounds(start, end),
	}
}

// LineText returns the raw bytes for the specified 1-based line number,
// excluding the trailing newline characters.
func (f *SourceText) LineText(number int) []byte {
	line := f.Line(number)
	if line.Empty() {
		return nil
	}

	return f.BytesOf(line.Span)
}

// searchLineIndex performs a fast O(log N) binary search to find the 0-based
// index of the line containing the given byte offset.
func (f *SourceText) searchLineIndex(offset int) int {
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
func (f *SourceText) LineAt(offset int) Line {
	return f.Line(f.searchLineIndex(offset) + 1)
}

// Positions

// Position translates a flat byte offset into a structural Position
// containing the 1-based line, byte column, and LSP-compliant column.
func (f *SourceText) Position(offset int) Position {
	offset = f.Clamp(offset)

	lineIdx := f.searchLineIndex(offset)
	lineStart := f.lineOffsets[lineIdx]
	byteColumn := offset - lineStart

	var lspColumn int

	prefixText := f.text[lineStart:offset]
	for len(prefixText) > 0 {
		r, size := utf8.DecodeRune(prefixText)

		if r >= 0x10000 {
			lspColumn += 2
		} else {
			lspColumn++
		}

		prefixText = prefixText[size:]
	}

	return Position{
		Offset:     offset,
		Line:       lineIdx + 1,
		ByteColumn: byteColumn + 1,
		LSPColumn:  lspColumn,
	}
}

// Offset converts a 1-based line and 1-based column into a flat byte offset.
func (f *SourceText) Offset(line, column int) int {
	l := f.Line(line)

	if l.Number == 0 {
		return f.Len()
	}

	if column <= 0 {
		return l.Span.Start()
	}

	return l.Span.Clamp(l.Span.Start() + column - 1)
}

// Spans

// BytesOf returns a slice of the file's text bounded by the given Span.
func (f *SourceText) BytesOf(span Span) []byte {
	sStart, sEnd := span.AsTuple()

	start := f.Clamp(sStart)
	end := max(f.Clamp(sEnd), start)

	return f.text[start:end]
}

// Text returns the string representation of the text bounded by the given Span.
func (f *SourceText) Text(span Span) string {
	return string(f.BytesOf(span))
}

// LineSpan returns a Span covering the entirety of the specified 1-based line.
func (f *SourceText) LineSpan(number int) Span {
	return f.Line(number).Span
}

// UTF-8

// Rune decodes and returns the rune at the given offset, along with its size in bytes.
func (f *SourceText) Rune(offset int) (r rune, size int) {
	offset = f.Clamp(offset)

	if offset >= f.Len() {
		return utf8.RuneError, 0
	}

	return utf8.DecodeRune(f.text[offset:])
}

// RuneBefore decodes and returns the rune immediately preceding the given offset.
func (f *SourceText) RuneBefore(offset int) (r rune, size int) {
	offset = f.Clamp(offset)

	if offset == 0 {
		return utf8.RuneError, 0
	}

	return utf8.DecodeLastRune(f.text[:offset])
}
