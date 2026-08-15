package text

import "unicode/utf8"

// Reader provides a highly optimized, forward-moving rune iterator.
// It is designed specifically to feed characters into the compiler's lexer.
type Reader struct {
	file   *SourceText
	text   []byte
	offset int
}

// NewReader initializes a character reader bound to the given file.
func NewReader(file *SourceText) *Reader {
	return &Reader{
		file: file,
		text: file.text,
	}
}

// Offset returns the reader's current byte position within the file.
func (r *Reader) Offset() int { return r.offset }

// EOF reports whether the reader has consumed all available bytes.
func (r *Reader) EOF() bool { return r.offset >= len(r.text) }

// Checkpoint saves the current reader position for arbitrary lookahead.
func (r *Reader) Checkpoint() int {
	return r.offset
}

// Restore rewinds the reader to a previously saved checkpoint.
// This replaces `Backup()` to allow infinite backtracking by the lexer.
func (r *Reader) Restore(offset int) {
	r.offset = offset
}

// Peek inspects the next rune and its byte size without advancing the reader.
// It favors an optimized path for standard ASCII characters to minimize decoding overhead.
func (r *Reader) Peek() (ch rune, size int) {
	if r.EOF() {
		return 0, 0
	}

	b := r.text[r.offset]
	if b < utf8.RuneSelf {
		return rune(b), 1 // Fast path for ASCII
	}

	return utf8.DecodeRune(r.text[r.offset:])
}

// Next consumes the next rune from the file and advances the reader state.
func (r *Reader) Next() (ch rune, size int) {
	ch, size = r.Peek()
	r.offset += size
	return ch, size
}
