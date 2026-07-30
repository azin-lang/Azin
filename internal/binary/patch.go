package binary

import (
	"encoding/binary"
	"math"
	"sync"
)

// patchPool recycles Emitter instances to prevent heap allocations
// when temporarily generating subcomponents during a Patch callback.
var patchPool = sync.Pool{
	New: func() any { return New(binary.NativeEndian) },
}

// PatchU8 overwrites a single byte at a historical offset.
func (e *Emitter) PatchU8(off uint64, value uint8) {
	e.data[off] = value
}

// PatchU16 overwrites a 16-bit integer at a historical offset.
func (e *Emitter) PatchU16(off uint64, value uint16) {
	e.order.PutUint16(e.data[off:], value)
}

// PatchU32 overwrites a 32-bit integer at a historical offset.
func (e *Emitter) PatchU32(off uint64, value uint32) {
	e.order.PutUint32(e.data[off:], value)
}

// PatchU64 overwrites a 64-bit integer at a historical offset.
func (e *Emitter) PatchU64(off uint64, value uint64) {
	e.order.PutUint64(e.data[off:], value)
}

// PatchI8 overwrites a signed 8-bit integer at a historical offset.
func (e *Emitter) PatchI8(off uint64, value int8) {
	e.PatchU8(off, uint8(value))
}

// PatchI16 overwrites a signed 16-bit integer at a historical offset.
func (e *Emitter) PatchI16(off uint64, value int16) {
	e.PatchU16(off, uint16(value))
}

// PatchI32 overwrites a signed 32-bit integer at a historical offset.
func (e *Emitter) PatchI32(off uint64, value int32) {
	e.PatchU32(off, uint32(value))
}

// PatchI64 overwrites a signed 64-bit integer at a historical offset.
func (e *Emitter) PatchI64(off uint64, value int64) {
	e.PatchU64(off, uint64(value))
}

// PatchF32 overwrites a 32-bit float at a historical offset.
func (e *Emitter) PatchF32(off uint64, value float32) {
	e.PatchU32(off, math.Float32bits(value))
}

// PatchF64 overwrites a 64-bit float at a historical offset.
func (e *Emitter) PatchF64(off uint64, value float64) {
	e.PatchU64(off, math.Float64bits(value))
}

// PatchPtr overwrites a pointer (4 or 8 bytes depending on size) at a historical offset.
func (e *Emitter) PatchPtr(off uint64, size int, value uint64) {
	switch size {
	case 4:
		e.PatchU32(off, uint32(value))
	case 8:
		e.PatchU64(off, value)
	default:
		panic("unsupported pointer size")
	}
}

// PatchBytes overwrites a span of bytes at a historical offset.
// Warning: This does not validate length and will panic if it exceeds capacity.
func (e *Emitter) PatchBytes(off uint64, data []byte) {
	copy(e.data[off:], data)
}

// PatchLabel32 overwrites a 32-bit integer located at the address of the provided label.
func (e *Emitter) PatchLabel32(label Label, value uint32) {
	e.PatchU32(uint64(label), value)
}

// PatchLabel64 overwrites a 64-bit integer located at the address of the provided label.
func (e *Emitter) PatchLabel64(label Label, value uint64) {
	e.PatchU64(uint64(label), value)
}

// PatchBool overwrites a boolean byte (1 or 0) at a historical offset.
func (e *Emitter) PatchBool(off uint64, value bool) {
	if value {
		e.PatchU8(off, 1)
	} else {
		e.PatchU8(off, 0)
	}
}

// PatchString overwrites bytes at a historical offset with the provided string.
func (e *Emitter) PatchString(off uint64, s string) {
	copy(e.data[off:], s)
}

// PatchCString overwrites bytes with a null-terminated C-style string.
func (e *Emitter) PatchCString(off uint64, s string) {
	copy(e.data[off:], s)
	e.data[off+uint64(len(s))] = 0
}

// Patch yields a temporary Emitter to a callback, allowing complex data structures
// to be built up before being copied into a historical offset in one block.
// This is heavily used to patch header structures where the sizes were not known upfront.
func (e *Emitter) Patch(off uint64, fn func(*Emitter)) {
	tmp := patchPool.Get().(*Emitter)

	// Reset length but keep backing array capacity
	tmp.Reset()
	tmp.order = e.order

	fn(tmp)

	copy(e.data[off:], tmp.data)

	patchPool.Put(tmp)
}

// AddU32 increments an existing 32-bit integer at 'off' by 'value'.
func (e *Emitter) AddU32(off uint64, value uint32) {
	v := e.order.Uint32(e.data[off:])
	e.PatchU32(off, v+value)
}

// AddU64 increments an existing 64-bit integer at 'off' by 'value'.
func (e *Emitter) AddU64(off uint64, value uint64) {
	v := e.order.Uint64(e.data[off:])
	e.PatchU64(off, v+value)
}

// PatchRelative32 patches a 32-bit jump offset at 'from' pointing to 'to',
// accounting for a standard 4-byte instruction width.
func (e *Emitter) PatchRelative32(from, to uint64) {
	disp := int32(to) - int32(from) - 4
	e.PatchI32(from, disp)
}

// AddU16 increments an existing 16-bit integer at 'off' by 'value'.
func (e *Emitter) AddU16(off uint64, value uint16) {
	v := e.order.Uint16(e.data[off:])
	e.PatchU16(off, v+value)
}

// AddI32 increments an existing signed 32-bit integer at 'off' by 'value'.
func (e *Emitter) AddI32(off uint64, value int32) {
	v := int32(e.order.Uint32(e.data[off:]))
	e.PatchI32(off, v+value)
}

// AddI64 increments an existing signed 64-bit integer at 'off' by 'value'.
func (e *Emitter) AddI64(off uint64, value int64) {
	v := int64(e.order.Uint64(e.data[off:]))
	e.PatchI64(off, v+value)
}

// PatchRelative64 patches a 64-bit jump offset at 'from' pointing to 'to',
// accounting for an 8-byte instruction width.
func (e *Emitter) PatchRelative64(from, to uint64) {
	e.PatchI64(from, int64(to)-int64(from)-8)
}

// PatchRelative16 patches a 16-bit jump offset at 'from' pointing to 'to',
// accounting for a 2-byte instruction width.
func (e *Emitter) PatchRelative16(from, to uint64) {
	e.PatchI16(from, int16(to)-int16(from)-2)
}

// CanPatch checks if writing 'size' bytes starting at 'off' fits within the
// currently generated data bounds.
func (e *Emitter) CanPatch(off, size uint64) bool {
	return off+size <= uint64(len(e.data))
}
