package text

// TextSpan represents a range of characters in the source text.
type TextSpan struct {
	Start  uint32
	Length uint32
}

func NewTextSpan(start, length uint32) TextSpan {
	return TextSpan{Start: start, Length: length}
}

func FromBounds(start, end uint32) TextSpan {
	return TextSpan{
		Start:  start,
		Length: end - start,
	}
}

func (s TextSpan) End() uint32 {
	return s.Start + s.Length
}

func (s TextSpan) IsEmpty() bool {
	return s.Length == 0
}

func (s TextSpan) Contains(position uint32) bool {
	return position >= s.Start && position < s.End()
}

func (s TextSpan) ContainsSpan(other TextSpan) bool {
	return other.Start >= s.Start && other.End() <= s.End()
}

func (s TextSpan) Overlaps(other TextSpan) bool {
	return s.Start < other.End() && other.Start < s.End()
}

func (s TextSpan) Intersects(other TextSpan) bool {
	return s.Overlaps(other)
}

func (s TextSpan) Intersection(other TextSpan) (TextSpan, bool) {
	start := max(s.Start, other.Start)
	end := min(s.End(), other.End())

	if start >= end {
		return TextSpan{}, false
	}

	return FromBounds(start, end), true
}

func (s TextSpan) Union(other TextSpan) TextSpan {
	start := min(s.Start, other.Start)
	end := max(s.End(), other.End())

	return FromBounds(start, end)
}

func (s TextSpan) WithStart(start uint32) TextSpan {
	return TextSpan{
		Start:  start,
		Length: s.Length,
	}
}

func (s TextSpan) WithLength(length uint32) TextSpan {
	return TextSpan{
		Start:  s.Start,
		Length: length,
	}
}

func (s TextSpan) Offset(offset int32) TextSpan {
	start := int64(s.Start) + int64(offset)

	return TextSpan{
		Start:  uint32(start),
		Length: s.Length,
	}
}

func (s TextSpan) Expand(amount uint32) TextSpan {
	return TextSpan{
		Start:  s.Start,
		Length: s.Length + amount,
	}
}

func (s TextSpan) TrimStart(amount uint32) TextSpan {
	if amount >= s.Length {
		return TextSpan{
			Start: s.End(),
		}
	}

	return TextSpan{
		Start:  s.Start + amount,
		Length: s.Length - amount,
	}
}

func (s TextSpan) TrimEnd(amount uint32) TextSpan {
	if amount >= s.Length {
		return TextSpan{
			Start: s.Start,
		}
	}

	return TextSpan{
		Start:  s.Start,
		Length: s.Length - amount,
	}
}

func (s TextSpan) Slice(text string) string {
	return text[s.Start:s.End()]
}

func (s TextSpan) Equals(other TextSpan) bool {
	return s.Start == other.Start && s.Length == other.Length
}
