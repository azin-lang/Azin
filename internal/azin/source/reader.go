package source

import "unicode/utf8"

// Reader provides a highly optimized, forward-moving rune iterator.
// It is designed specifically to feed characters into the compiler's lexer.
type Reader struct {
	file       *File
	text       []byte
	offset     uint32
	prevOffset uint32
}

// NewReader initializes a character reader bound to the given file.
func NewReader(file *File) *Reader {
	return &Reader{
		file: file,
		text: file.text,
	}
}

// Offset returns the reader's current byte position within the file.
func (r *Reader) Offset() uint32 { return r.offset }

// EOF reports whether the reader has consumed all available bytes.
func (r *Reader) EOF() bool { return r.offset >= uint32(len(r.text)) }

// Peek inspects the next rune and its byte size without advancing the reader.
// It favors an optimized path for standard ASCII characters to minimize decoding overhead.
func (r *Reader) Peek() (rune, uint32) {
	if r.EOF() {
		return 0, 0
	}

	b := r.text[r.offset]
	if b < utf8.RuneSelf {
		return rune(b), 1
	}

	ru, sz := utf8.DecodeRune(r.text[r.offset:])
	return ru, uint32(sz)
}

// Next consumes the next rune from the file and advances the reader state.
func (r *Reader) Next() (rune, uint32) {
	r.prevOffset = r.offset
	ch, size := r.Peek()
	r.offset += size
	return ch, size
}

// Backup retreats the reader by exactly one rune, reverting to the previous state.
func (r *Reader) Backup() {
	r.offset = r.prevOffset
}
