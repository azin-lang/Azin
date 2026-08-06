package text

// Line represents the boundaries of a single, physical line of text in a file.
// The bounds cover only the visible characters, expressly excluding \r or \n.
type Line struct {
	// Number is the 1-based line sequence number.
	Number uint32

	// Start is the byte offset of the first character in the line.
	Start uint32

	// End is the byte offset immediately following the last character in the line
	// (before any line-ending characters).
	End uint32
}

// Len returns the total byte length of the printable characters on this line.
func (l Line) Len() uint32 { return l.End - l.Start }

// Empty reports whether the line contains zero printable bytes.
func (l Line) Empty() bool { return l.Start == l.End }
