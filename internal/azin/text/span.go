package text

import "fmt"

// Span represents a half-open byte range [Start, End) within a file.
// Start is inclusive, pointing to the first byte of the range.
// End is exclusive, pointing to the byte immediately following the range.
type Span struct {
	Start uint32
	End   uint32
}

// NewSpan constructs a new half-open span [start, end).
func NewSpan(start, end uint32) Span {
	return Span{Start: start, End: end}
}

// NewPointSpan creates a zero-length span at a specific byte offset.
// This is useful for representing single-character tokens or insertion points.
func NewPointSpan(offset uint32) Span {
	return Span{Start: offset, End: offset}
}

// Valid reports whether the span has a logical layout (Start is before or equal to End).
func (s Span) Valid() bool { return s.Start <= s.End }

// Empty reports whether the span has a length of zero.
func (s Span) Empty() bool { return s.Start == s.End }

// Len returns the number of bytes enclosed by the span.
func (s Span) Len() uint32 { return s.End - s.Start }

// Intersection calculates the overlapping span between this span and another.
// It returns an empty span at the nearest boundary if no overlap exists.
func (s Span) Intersection(other Span) Span {
	start := max(s.Start, other.Start)
	end := max(min(s.End, other.End), start)
	return NewSpan(start, end)
}

// Cover creates the smallest bounding span that completely encloses both spans.
// Used heavily by the parser to merge child expression spans into parent AST node spans.
func (s Span) Cover(other Span) Span {
	return NewSpan(min(s.Start, other.Start), max(s.End, other.End))
}

// String provides a debugging string representation of the span.
func (s Span) String() string {
	return fmt.Sprintf("%d..%d", s.Start, s.End)
}

// IsValidAndNonEmpty is a convenience method ensuring the span is logically sound and holds data.
func (s Span) IsValidAndNonEmpty() bool {
	return s.Valid() && !s.Empty()
}
