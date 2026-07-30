package binary

import (
	"encoding/binary"
	"unsafe"
)

// U32Slice writes a slice of 32-bit integers to the emitter.
func (e *Emitter) U32Slice(v []uint32) {
	if len(v) == 0 {
		return
	}

	// FAST PATH: If target endianness matches the host architecture,
	// bypass encoding entirely and use an unsafe direct memory copy.
	if e.order == binary.NativeEndian {
		b := unsafe.Slice((*byte)(unsafe.Pointer(&v[0])), len(v)*4)
		e.data = append(e.data, b...)
		return
	}

	start := len(e.data)
	e.data = append(e.data, make([]byte, len(v)*4)...)
	buf := e.data[start:]
	for i, val := range v {
		e.order.PutUint32(buf[i*4:], val)
	}
}

// U64Slice writes a slice of 64-bit integers to the emitter.
func (e *Emitter) U64Slice(v []uint64) {
	if len(v) == 0 {
		return
	}

	// FAST PATH: If target endianness matches the host architecture,
	// bypass encoding entirely and use an unsafe direct memory copy.
	if e.order == binary.NativeEndian {
		b := unsafe.Slice((*byte)(unsafe.Pointer(&v[0])), len(v)*8)
		e.data = append(e.data, b...)
		return
	}

	start := len(e.data)
	e.data = append(e.data, make([]byte, len(v)*8)...)
	buf := e.data[start:]
	for i, val := range v {
		e.order.PutUint64(buf[i*8:], val)
	}
}
