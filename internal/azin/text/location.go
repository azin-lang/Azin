package text

import "fmt"

// Location combines a physical file pointer with a specific byte span.
// It is embedded in every AST node to attach syntax trees to their origin text.
type Location struct {
	File *SourceText
	Span Span
}

// Position resolves the starting span boundary into human-readable coordinates.
func (l Location) Position() Position {
	if l.File == nil {
		return Position{}
	}

	return l.File.Position(l.Span.Start)
}

// String returns a formatted identifier denoting the file and exact coordinates
// (e.g., "main.az:12:4").
func (l Location) String() string {
	if l.File == nil {
		return l.Span.String()
	}

	return fmt.Sprintf("%s:%s", l.File.Base(), l.Position())
}
