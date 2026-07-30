package binary

// ULEB128 encodes an unsigned 64-bit integer using Little-Endian Base 128 compression.
// Often used in DWARF debugging info and WebAssembly to save space on small integers.
func (e *Emitter) ULEB128(v uint64) {
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80 // Set continuation bit
		}
		e.U8(c)
		if v == 0 {
			break
		}
	}
}

// SLEB128 encodes a signed 64-bit integer using Little-Endian Base 128 compression.
func (e *Emitter) SLEB128(v int64) {
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		signBit := (c & 0x40) != 0

		// If the value is 0 and the sign bit is not set, or
		// if the value is -1 and the sign bit is set, we are done.
		if (v == 0 && !signBit) || (v == -1 && signBit) {
			e.U8(c)
			break
		}

		e.U8(c | 0x80) // Set continuation bit
	}
}
