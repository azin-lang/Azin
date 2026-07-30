package binary

// Align calculates the next offset that satisfies the given alignment boundary.
// If the alignment is 1 or less, it returns the offset unchanged.
func Align(offset, alignment uint64) uint64 {
	if alignment <= 1 {
		return offset
	}
	return (offset + alignment - 1) &^ (alignment - 1)
}

// Padding calculates the number of zero bytes required to reach the next alignment boundary.
func Padding(offset, alignment uint64) uint64 {
	return Align(offset, alignment) - offset
}

// IsAligned returns true if the current offset naturally falls on the alignment boundary.
func IsAligned(offset, alignment uint64) bool {
	return Padding(offset, alignment) == 0
}

// Align pads the emitter with zeroes until the current offset matches the given alignment.
func (e *Emitter) Align(alignment uint64) {
	e.Pad(Padding(e.Offset(), alignment))
}

// Pad appends 'count' zeroes to the emitter.
func (e *Emitter) Pad(count uint64) {
	e.Fill(0, count)
}

// Skip advances the emitter by 'count' bytes, filling the skipped space with zeroes.
func (e *Emitter) Skip(count uint64) {
	e.Pad(count)
}

// AlignDown calculates the nearest preceding offset that satisfies the alignment boundary.
func AlignDown(offset, alignment uint64) uint64 {
	if alignment <= 1 {
		return offset
	}
	return offset &^ (alignment - 1)
}

// NextAlignment is a convenience alias for Padding, returning the bytes to the next boundary.
func NextAlignment(offset, alignment uint64) uint64 {
	return Padding(offset, alignment)
}

// AlignCode aligns the emitter to a 16-byte boundary, typically optimal for instruction caching.
func (e *Emitter) AlignCode() {
	e.Align(16)
}

// AlignData aligns the emitter to an 8-byte boundary, standard for 64-bit data structures.
func (e *Emitter) AlignData() {
	e.Align(8)
}

// AlignPointer aligns the emitter to match the target architecture's pointer size.
func (e *Emitter) AlignPointer(ptrSize int) {
	e.Align(uint64(ptrSize))
}

// AlignFill pads the emitter with a specific byte pattern until the alignment boundary is reached.
// This is useful for filling code alignment gaps with NOP instructions (e.g., 0x90).
func (e *Emitter) AlignFill(alignment uint64, fill byte) {
	e.Fill(fill, Padding(e.Offset(), alignment))
}
