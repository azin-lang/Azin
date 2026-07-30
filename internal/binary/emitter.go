// ===== FILE: internal/binary/emitter.go =====
package binary

import (
	"encoding/binary"
	"slices"
)

// ByteOrder combines the classic binary interface (for patching memory)
// with the modern AppendByteOrder interface (for zero-allocation slice extension).
type ByteOrder interface {
	binary.ByteOrder
	binary.AppendByteOrder
}

// Emitter is a high-performance buffer for generating and patching binary data.
// It is designed for low-overhead assembly generation and executable formatting.
type Emitter struct {
	data  []byte
	order ByteOrder
}

// New creates an Emitter initialized with the target endianness.
func New(order ByteOrder) *Emitter {
	return &Emitter{
		order: order,
	}
}

// Data returns the underlying byte slice representing the emitted binary.
func (e *Emitter) Data() []byte { return e.data }

// Order returns the configured endianness of the emitter.
func (e *Emitter) Order() ByteOrder { return e.order }

// Cap returns the total allocated capacity of the underlying slice.
func (e *Emitter) Cap() int { return cap(e.data) }

// Len returns the current length (in bytes) of the emitted data.
func (e *Emitter) Len() int { return len(e.data) }

// Offset returns the current write cursor position (same as Len).
func (e *Emitter) Offset() uint64 { return uint64(e.Len()) }

// Size returns the total size in bytes (same as Offset).
func (e *Emitter) Size() uint64 { return uint64(e.Len()) }

// Empty returns true if no data has been emitted yet.
func (e *Emitter) Empty() bool { return e.Len() == 0 }

// Reset clears the emitted data by setting the slice length to zero.
// It retains the underlying memory capacity for future writes.
func (e *Emitter) Reset() {
	e.data = e.data[:0]
}

// Grow ensures that at least 'n' more bytes can be written without triggering an allocation.
func (e *Emitter) Grow(n uint64) {
	if e.CapacityLeft() < int(n) {
		e.data = slices.Grow(e.data, int(n))
	}
}

// Clone returns a deep copy of the emitted binary data.
func (e *Emitter) Clone() []byte {
	return slices.Clone(e.data)
}

// Truncate rolls back the emitter's length to the specified offset, discarding data written after it.
func (e *Emitter) Truncate(offset uint64) {
	e.data = e.data[:offset]
}

// SetOffset manually adjusts the write cursor. If the target offset is greater than
// the current length, it pads the gap with zeroes. If it's smaller, it truncates.
func (e *Emitter) SetOffset(offset uint64) {
	if offset <= uint64(e.Len()) {
		e.data = e.data[:offset]
		return
	}
	e.Pad(offset - uint64(e.Len()))
}

// Clear completely resets the emitter.
func (e *Emitter) Clear() {
	e.Reset()
}

// Release nils the internal buffer, allowing the garbage collector to free the memory.
func (e *Emitter) Release() {
	e.data = nil
}

// CapacityLeft returns how many bytes can be added before the underlying slice must reallocate.
func (e *Emitter) CapacityLeft() int {
	return cap(e.data) - len(e.data)
}

// ReserveCapacity forces the underlying slice to grow to accommodate 'n' total bytes.
func (e *Emitter) ReserveCapacity(n int) {
	e.data = slices.Grow(e.data, n)
}

// Compact reallocates the buffer to exactly fit the data, freeing unused trailing capacity.
func (e *Emitter) Compact() {
	e.data = slices.Clip(e.data)
}

// LastOffset returns the byte offset immediately preceding the current write cursor.
func (e *Emitter) LastOffset() uint64 {
	if len(e.data) == 0 {
		return 0
	}
	return uint64(len(e.data) - 1)
}

// RemainingCapacity is an alias for CapacityLeft.
func (e *Emitter) RemainingCapacity() int {
	return cap(e.data) - len(e.data)
}
