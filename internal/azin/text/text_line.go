package text

// Line represents the boundaries of a single, physical line of text in a file.
// The bounds cover only the visible characters, expressly excluding \r or \n.
type Line struct {
	Number int
	Span   Span
}

// Len returns the total byte length of the printable characters on this line.
func (l Line) Len() int { return l.Span.Len() }

// Empty reports whether the line contains zero printable bytes.
func (l Line) Empty() bool { return l.Span.IsEmpty() }
