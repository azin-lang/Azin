package syntax

import (
	"github.com/azin-lang/Azin/internal/azin/source"
)

// Trivia represents formatting elements like whitespace or comments.
type Trivia struct {
	Kind Kind
	Span source.Span
}

// Token represents a code terminal, owning its layout via trivia.
type Token struct {
	Kind           Kind
	Span           source.Span // Position of the token text itself
	LeadingTrivia  []Trivia    // Trivia preceding the token
	TrailingTrivia []Trivia    // Trivia following the token on the same line
}

// FullSpan returns the absolute span including leading and trailing trivia.
func (t Token) FullSpan() source.Span {
	start := t.Span.Start
	if len(t.LeadingTrivia) > 0 {
		start = t.LeadingTrivia[0].Span.Start
	}

	end := t.Span.End
	if len(t.TrailingTrivia) > 0 {
		end = t.TrailingTrivia[len(t.TrailingTrivia)-1].Span.End
	}

	return source.NewSpan(start, end)
}

// Text returns the exact source text corresponding to the token itself.
func (t Token) Text(file *source.File) string {
	return file.Text(t.Span)
}

// FullText returns the text of the token including all attached trivia.
func (t Token) FullText(file *source.File) string {
	return file.Text(t.FullSpan())
}

// HasLeadingTrivia reports whether the token has any leading trivia.
func (t Token) HasLeadingTrivia() bool {
	return len(t.LeadingTrivia) > 0
}

// HasTrailingTrivia reports whether the token has any trailing trivia.
func (t Token) HasTrailingTrivia() bool {
	return len(t.TrailingTrivia) > 0
}
