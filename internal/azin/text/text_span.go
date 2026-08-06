package text

import "fmt"

// Span represents an immutable half-open byte interval [start, end).
//
// A Span identifies a contiguous region within a source file or buffer.
// It owns no memory, is highly efficient to copy, and guarantees:
//
//	0 <= start <= end <= MaxUint32
type Span struct {
	start uint32
	end   uint32
}

const maxUint32 = ^uint32(0)

func addClamp(a, b uint32) uint32 {
	if maxUint32-a < b {
		return maxUint32
	}
	return a + b
}

func subClamp(a, b uint32) uint32 {
	if a < b {
		return 0
	}
	return a - b
}

//
// Constructors
//

// NewSpan constructs a span from a start offset and length.
// Safely clamps to MaxUint32 to prevent overflow.
func NewSpan(start, length uint32) Span {
	if length == 0 {
		return SpanAt(start)
	}
	return Span{
		start: start,
		end:   addClamp(start, length),
	}
}

// SpanFromBounds constructs a span from explicit start and end markers.
// If end precedes start, it normalizes by collapsing to an empty span at start.
func SpanFromBounds(start, end uint32) Span {
	return Span{
		start: start,
		end:   max(start, end),
	}
}

// SpanAt constructs an empty span (a cursor) located exactly at offset.
func SpanAt(offset uint32) Span {
	return Span{
		start: offset,
		end:   offset,
	}
}

//
// Properties
//

// Start returns the first byte offset.
func (s Span) Start() uint32 {
	return s.start
}

// End returns the first byte immediately after the span.
func (s Span) End() uint32 {
	return s.end
}

// Len returns the span length in bytes.
func (s Span) Len() uint32 {
	return s.end - s.start
}

// IsEmpty reports whether the span contains no bytes.
func (s Span) IsEmpty() bool {
	return s.start == s.end
}

// AsTuple returns both the start and end offsets as a raw pair.
func (s Span) AsTuple() (start, end uint32) {
	return s.start, s.end
}

//
// Relationships
//

// Equals reports whether both spans have identical bounds.
func (s Span) Equals(other Span) bool {
	return s == other
}

// ContainsOffset reports whether offset lies strictly inside the span.
// Optimized with an early exit if the span is empty.
func (s Span) ContainsOffset(offset uint32) bool {
	if s.IsEmpty() {
		return false
	}
	return offset >= s.start && offset < s.end
}

// Encloses reports whether the other span lies entirely inside this span.
func (s Span) Encloses(other Span) bool {
	return other.start >= s.start && other.end <= s.end
}

// Overlaps reports whether the spans share at least one byte.
// Empty spans contain no bytes and therefore never overlap anything.
func (s Span) Overlaps(other Span) bool {
	if s.IsEmpty() || other.IsEmpty() {
		return false
	}
	return s.start < other.end && other.start < s.end
}

// Touches reports whether the spans meet without overlapping.
func (s Span) Touches(other Span) bool {
	return s.end == other.start || other.end == s.start
}

// Distance returns the gap in bytes between two non-overlapping spans.
// Returns 0 if they touch, overlap, or enclose each other.
func (s Span) Distance(other Span) uint32 {
	if s.Overlaps(other) || s.Encloses(other) || other.Encloses(s) || s.Touches(other) {
		return 0
	}
	if s.IsBefore(other) {
		return other.start - s.end
	}
	return s.start - other.end
}

// IsBefore reports whether this span ends before or exactly where other begins.
func (s Span) IsBefore(other Span) bool {
	return s.end <= other.start
}

// IsAfter reports whether this span begins after or exactly where other ends.
func (s Span) IsAfter(other Span) bool {
	return s.start >= other.end
}

//
// Transformations
//

// Intersection returns the overlapping region of two spans.
// It returns false if the spans are completely disjoint.
func (s Span) Intersection(other Span) (Span, bool) {
	if s == other {
		return s, true
	}
	if s.end <= other.start || other.end <= s.start {
		return Span{}, false
	}

	start := max(s.start, other.start)
	end := min(s.end, other.end)

	if start < end || (start == end && (s.Encloses(other) || other.Encloses(s))) {
		return SpanFromBounds(start, end), true
	}

	return Span{}, false
}

// Encompass returns the smallest possible bounding span that completely
// covers both spans, including any empty space between them.
// Optimized with container checks to bypass calculation overhead.
func (s Span) Encompass(other Span) Span {
	if s.Encloses(other) {
		return s
	}
	if other.Encloses(s) {
		return other
	}
	return Span{
		start: min(s.start, other.start),
		end:   max(s.end, other.end),
	}
}

// Expand grows the span outward by left and right bytes.
// Safely clamps at boundaries. Early exit if arguments are zero.
func (s Span) Expand(left, right uint32) Span {
	if left == 0 && right == 0 {
		return s
	}
	return Span{
		start: subClamp(s.start, left),
		end:   addClamp(s.end, right),
	}
}

// Shrink removes bytes from each side.
// If the requested shrink would invert the span, an empty span is returned.
func (s Span) Shrink(left, right uint32) Span {
	if left == 0 && right == 0 {
		return s
	}

	length := s.Len()
	if left >= length || right >= length-left {
		safeStart := addClamp(s.start, min(left, length))
		return SpanAt(safeStart)
	}

	return Span{
		start: s.start + left,
		end:   s.end - right,
	}
}

// Translate shifts the entire span by a delta amount.
// Accepts int64 to safely handle massive shifts without bit truncation.
func (s Span) Translate(delta int64) Span {
	if delta == 0 {
		return s
	}

	if delta > 0 {
		if delta >= int64(maxUint32) {
			return Span{start: maxUint32, end: maxUint32}
		}
		d := uint32(delta)
		return Span{
			start: addClamp(s.start, d),
			end:   addClamp(s.end, d),
		}
	}

	if delta <= -int64(maxUint32) {
		return Span{start: 0, end: 0}
	}

	d := uint32(-delta)
	return Span{
		start: subClamp(s.start, d),
		end:   subClamp(s.end, d),
	}
}

// Clamp constrains an offset to lie within the boundaries of the span.
// Optimized with fast-path bounds checks.
func (s Span) Clamp(offset uint32) uint32 {
	if offset <= s.start {
		return s.start
	}
	if offset >= s.end {
		return s.end
	}
	return offset
}

// SplitAt divides the span into two adjacent spans at the given offset.
// Offsets outside the span are safely clamped.
func (s Span) SplitAt(offset uint32) (Span, Span) {
	offset = s.Clamp(offset)
	return Span{start: s.start, end: offset}, Span{start: offset, end: s.end}
}

//
// Immutable modifiers
//

// WithStart returns a copy with a different start.
func (s Span) WithStart(start uint32) Span {
	if start == s.start {
		return s
	}
	return SpanFromBounds(start, s.end)
}

// WithEnd returns a copy with a different end.
func (s Span) WithEnd(end uint32) Span {
	if end == s.end {
		return s
	}
	return SpanFromBounds(s.start, end)
}

// WithLength returns a copy with a different length.
func (s Span) WithLength(length uint32) Span {
	if length == s.Len() {
		return s
	}
	return Span{
		start: s.start,
		end:   addClamp(s.start, length),
	}
}

//
// Formatting
//

// String returns the span mathematically formatted as a half-open interval.
func (s Span) String() string {
	return fmt.Sprintf("[%d, %d)", s.start, s.end)
}
