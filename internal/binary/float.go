package binary

import "math"

// F32 converts a 32-bit float to IEEE 754 bits and emits it.
func (e *Emitter) F32(v float32) {
	e.U32(math.Float32bits(v))
}

// F64 converts a 64-bit float to IEEE 754 bits and emits it.
func (e *Emitter) F64(v float64) {
	e.U64(math.Float64bits(v))
}

// F80 emits an 80-bit extended-precision float (commonly used in x87 FPU architecture).
func (e *Emitter) F80(data [10]byte) {
	e.Bytes(data[:])
}
