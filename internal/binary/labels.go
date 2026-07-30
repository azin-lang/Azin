package binary

// Label represents a symbolic offset within the emitter.
// It is used to record locations for future backpatching (e.g., jump targets).
type Label uint64

// Label marks the current write cursor position and returns it as a Label.
func (e *Emitter) Label() Label {
	return Label(e.Offset())
}

// DistanceTo calculates the number of bytes from a historical label to the current offset.
func (e *Emitter) DistanceTo(label Label) uint64 {
	return e.Offset() - uint64(label)
}

// OffsetOf retrieves the raw byte offset associated with a given Label.
func (e *Emitter) OffsetOf(label Label) uint64 {
	return uint64(label)
}

// DistanceBetween calculates the byte span separating two labels.
func (e *Emitter) DistanceBetween(a, b Label) uint64 {
	return e.OffsetOf(b) - e.OffsetOf(a)
}
