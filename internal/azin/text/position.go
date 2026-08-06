package text

import "fmt"

// Position represents a human-readable and LSP-compliant coordinate within a file.
type Position struct {
	// Offset is the flat, 0-based byte index from the start of the file.
	Offset uint32

	// Line is the 1-based line number.
	Line uint32

	// ByteColumn is the 1-based byte offset from the start of the line.
	// Used primarily for formatting CLI terminal errors (e.g. main.az:10:5).
	ByteColumn uint32

	// LSPColumn is the 0-based UTF-16 code unit offset from the start of the line.
	// This is strictly required when formatting data for Language Server Protocol clients.
	LSPColumn uint32
}

// Valid reports whether the position contains a valid line number (greater than 0).
func (p Position) Valid() bool {
	return p.Line != 0
}

// String formats the position in a standard "line:column" notation.
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.ByteColumn)
}
