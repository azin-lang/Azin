package binary

// Reserve advances the emitter by 'size' bytes, leaving zeroes, and returns
// the original offset. This is typically used to reserve space for headers
// or sizes that will be backpatched later.
func (e *Emitter) Reserve(size uint64) uint64 {
	off := e.Offset()
	e.Pad(size)
	return off
}

// ReserveU8 reserves 1 byte for future patching.
func (e *Emitter) ReserveU8() uint64 {
	return e.Reserve(1)
}

// ReserveU16 reserves 2 bytes for future patching.
func (e *Emitter) ReserveU16() uint64 {
	return e.Reserve(2)
}

// ReserveU32 reserves 4 bytes for future patching.
func (e *Emitter) ReserveU32() uint64 {
	return e.Reserve(4)
}

// ReserveU64 reserves 8 bytes for future patching.
func (e *Emitter) ReserveU64() uint64 {
	return e.Reserve(8)
}

// Distance calculates the delta in bytes from the specified offset to the current tip.
func (e *Emitter) Distance(off uint64) uint64 {
	return e.Offset() - off
}

// ReservePtr reserves the target architecture's pointer size (4 or 8 bytes) for future patching.
func (e *Emitter) ReservePtr(ptrSize int) uint64 {
	return e.Reserve(uint64(ptrSize))
}

// ReserveCString reserves 'length' bytes plus 1 (for the null terminator).
func (e *Emitter) ReserveCString(length uint64) uint64 {
	return e.Reserve(length + 1)
}

// Advance forcefully moves the write cursor forward, bypassing normal writing methods.
func (e *Emitter) Advance(size uint64) {
	e.SetOffset(e.Offset() + size)
}

// ReserveAlignment aligns the emitter to the specified boundary, returning the offset
// at which the alignment actually started (post-padding).
func (e *Emitter) ReserveAlignment(alignment uint64) uint64 {
	off := e.Offset()
	e.Align(alignment)
	return off
}

// PlaceholderU8 emits a single zero byte and returns its offset.
func (e *Emitter) PlaceholderU8() uint64 {
	off := e.Offset()
	e.U8(0)
	return off
}

// PlaceholderU16 emits two zero bytes and returns their offset.
func (e *Emitter) PlaceholderU16() uint64 {
	off := e.Offset()
	e.U16(0)
	return off
}

// PlaceholderU32 emits four zero bytes and returns their offset.
func (e *Emitter) PlaceholderU32() uint64 {
	off := e.Offset()
	e.U32(0)
	return off
}

// PlaceholderU64 emits eight zero bytes and returns their offset.
func (e *Emitter) PlaceholderU64() uint64 {
	off := e.Offset()
	e.U64(0)
	return off
}

// PlaceholderPtr emits zeroed bytes equal to the pointer size and returns their offset.
func (e *Emitter) PlaceholderPtr(size int) uint64 {
	switch size {
	case 4:
		return e.PlaceholderU32()
	case 8:
		return e.PlaceholderU64()
	default:
		panic("unsupported pointer size")
	}
}
