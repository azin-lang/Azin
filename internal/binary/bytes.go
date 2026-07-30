package binary

// Byte appends a single 8-bit value to the emitter.
func (e *Emitter) Byte(v byte) {
	e.U8(v)
}

// Bool appends a boolean as a single byte (1 for true, 0 for false).
func (e *Emitter) Bool(v bool) {
	if v {
		e.Byte(1)
	} else {
		e.Byte(0)
	}
}

// Bytes appends a raw byte slice to the emitter.
func (e *Emitter) Bytes(data []byte) {
	e.data = append(e.data, data...)
}

// Fill appends 'count' copies of the given byte 'b' to the emitter.
// It uses pre-allocation and slice copying for high performance.
func (e *Emitter) Fill(b byte, count uint64) {
	if count == 0 {
		return
	}

	start := len(e.data)
	e.data = append(e.data, make([]byte, count)...)

	if b != 0 {
		buf := e.data[start:]
		for i := range buf {
			buf[i] = b
		}
	}
}

// RepeatU32 appends the given 32-bit integer 'count' times.
// Uses SIMD-friendly exponential doubling (memmove) to avoid loop overhead.
func (e *Emitter) RepeatU32(v uint32, count uint64) {
	if count == 0 {
		return
	}

	start := len(e.data)
	size := int(count * 4)
	e.Grow(size)
	e.data = e.data[:start+size]

	buf := e.data[start:]
	e.order.PutUint32(buf, v)

	for i := 4; i < size; {
		i += copy(buf[i:], buf[:i])
	}
}

// RepeatU64 appends the given 64-bit integer 'count' times.
// Uses SIMD-friendly exponential doubling for bulk emission speed.
func (e *Emitter) RepeatU64(v uint64, count uint64) {
	if count == 0 {
		return
	}

	start := len(e.data)
	size := int(count * 8)
	e.Grow(size)
	e.data = e.data[:start+size]

	buf := e.data[start:]
	e.order.PutUint64(buf, v)

	for i := 8; i < size; {
		i += copy(buf[i:], buf[:i])
	}
}

// Zeroes appends 'count' null bytes (0x00) to the emitter.
func (e *Emitter) Zeroes(count uint64) {
	e.Fill(0, count)
}

// Ones appends 'count' bytes of all 1s (0xFF) to the emitter.
func (e *Emitter) Ones(count uint64) {
	e.Fill(0xFF, count)
}

// RepeatBytes repeats the provided byte slice 'count' times.
// Optimized with exponential doubling to rapidly fill memory.
func (e *Emitter) RepeatBytes(data []byte, count uint64) {
	if count == 0 || len(data) == 0 {
		return
	}

	start := len(e.data)
	size := int(count) * len(data)
	e.Grow(size)
	e.data = e.data[:start+size]

	buf := e.data[start:]
	copy(buf, data)

	for i := len(data); i < size; {
		i += copy(buf[i:], buf[:i])
	}
}

// ByteAt reads a single byte at the historical offset.
func (e *Emitter) ByteAt(off uint64) byte {
	return e.data[off]
}

// Slice returns a sub-slice of the emitted data.
//
// Note: This returns a reference, not a copy. Modifying it modifies the emitted binary.
func (e *Emitter) Slice(off, size uint64) []byte {
	return e.data[off : off+size]
}

// ZeroUntil pads the emitter with zeroes until it reaches the specified absolute offset.
func (e *Emitter) ZeroUntil(offset uint64) {
	if offset <= e.Offset() {
		return
	}
	e.Zeroes(offset - e.Offset())
}

// Append adds the entire contents of another Emitter to this one.
func (e *Emitter) Append(other *Emitter) {
	e.data = append(e.data, other.data...)
}

// AppendAt copies the contents of another Emitter into a specific historical offset.
//
// Warning: This does not grow the buffer and will panic if it overruns the capacity.
func (e *Emitter) AppendAt(off uint64, other *Emitter) {
	copy(e.data[off:], other.data)
}
