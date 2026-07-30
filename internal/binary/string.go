package binary

import "unicode/utf8"

// String appends a standard Go string to the emitter without any length prefix or terminator.
func (e *Emitter) String(s string) {
	e.data = append(e.data, s...)
}

// CString appends a string followed by a null byte terminator (0x00), common in C APIs.
func (e *Emitter) CString(s string) {
	e.String(s)
	e.U8(0)
}

// Rune appends a single Unicode code point using UTF-8 encoding.
func (e *Emitter) Rune(r rune) {
	e.data = utf8.AppendRune(e.data, r)
}

// PascalString appends a string prefixed by its length as a single byte (max length 255).
func (e *Emitter) PascalString(s string) {
	if len(s) > 255 {
		panic("binary: Pascal string exceeds 255 bytes")
	}
	e.U8(uint8(len(s)))
	e.String(s)
}

// U16String appends a string prefixed by its length as a 16-bit integer (max length 65535).
func (e *Emitter) U16String(s string) {
	if len(s) > 65535 {
		panic("binary: U16 string exceeds 65535 bytes")
	}
	e.U16(uint16(len(s)))
	e.String(s)
}

// U32String appends a string prefixed by its length as a 32-bit integer.
func (e *Emitter) U32String(s string) {
	// no need for a check here, realistically who is doing a 4 gigabyte string?
	e.U32(uint32(len(s)))
	e.String(s)
}
