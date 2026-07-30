package binary

// U8 appends an unsigned 8-bit integer.
func (e *Emitter) U8(v uint8) {
	e.data = append(e.data, v)
}

// U16 appends an unsigned 16-bit integer using target endianness.
// Uses fast-path compiler intrinsic appending to avoid allocations.
func (e *Emitter) U16(v uint16) {
	e.data = e.order.AppendUint16(e.data, v)
}

// U32 appends an unsigned 32-bit integer using target endianness.
// Uses fast-path compiler intrinsic appending to avoid allocations.
func (e *Emitter) U32(v uint32) {
	e.data = e.order.AppendUint32(e.data, v)
}

// U64 appends an unsigned 64-bit integer using target endianness.
// Uses fast-path compiler intrinsic appending to avoid allocations.
func (e *Emitter) U64(v uint64) {
	e.data = e.order.AppendUint64(e.data, v)
}

// I8 appends a signed 8-bit integer.
func (e *Emitter) I8(v int8) { e.U8(uint8(v)) }

// I16 appends a signed 16-bit integer.
func (e *Emitter) I16(v int16) { e.U16(uint16(v)) }

// I32 appends a signed 32-bit integer.
func (e *Emitter) I32(v int32) { e.U32(uint32(v)) }

// I64 appends a signed 64-bit integer.
func (e *Emitter) I64(v int64) { e.U64(uint64(v)) }

// Uint dynamically sizes an unsigned int based on the host architecture
// (32-bit or 64-bit) and emits it.
func (e *Emitter) Uint(v uint) {
	if ^uint(0)>>32 == 0 {
		e.U32(uint32(v))
	} else {
		e.U64(uint64(v))
	}
}

// Int dynamically sizes a signed int based on the host architecture and emits it.
func (e *Emitter) Int(v int) {
	if ^uint(0)>>32 == 0 {
		e.I32(int32(v))
	} else {
		e.I64(int64(v))
	}
}

// Ptr emits a target-specific pointer size (either 4 or 8 bytes) as an unsigned integer.
func (e *Emitter) Ptr(size int, value uint64) {
	switch size {
	case 4:
		e.U32(uint32(value))
	case 8:
		e.U64(value)
	default:
		panic("unsupported pointer size")
	}
}

// IPtr emits a target-specific pointer size (either 4 or 8 bytes) as a signed integer.
func (e *Emitter) IPtr(size int, value int64) {
	switch size {
	case 4:
		e.I32(int32(value))
	case 8:
		e.I64(value)
	default:
		panic("unsupported pointer size")
	}
}
