package x86_64

import (
	stdbinary "encoding/binary"

	"github.com/azin-lang/Azin/internal/binary"
)

// Emitter emits raw x86-64 machine code.
type Emitter struct {
	*binary.Emitter
}

// New creates a new x86-64 emitter.
func New() *Emitter {
	return &Emitter{
		Emitter: binary.New(binary.ByteOrder(stdbinary.LittleEndian)),
	}
}

// MovEAX emits:
//
//	mov eax, imm32
func (e *Emitter) MovEAX(imm uint32) {
	e.U8(0xB8)
	e.U32(imm)
}

// MovEDI emits:
//
//	mov edi, imm32
func (e *Emitter) MovEDI(imm uint32) {
	e.U8(0xBF)
	e.U32(imm)
}

// Syscall emits:
//
//	syscall
func (e *Emitter) Syscall() {
	e.U16(0x050F) // Little-endian: 0F 05
}

// Ret emits:
//
//	ret
func (e *Emitter) Ret() {
	e.U8(0xC3)
}

// Nop emits:
//
//	nop
func (e *Emitter) Nop() {
	e.U8(0x90)
}

// EmitMinimalExit returns code equivalent to:
//
//	mov eax, 60
//	mov edi, 42
//	syscall
func EmitMinimalExit() *Emitter {
	e := New()

	e.MovEAX(60) // sys_exit
	e.MovEDI(42) // exit code
	e.Syscall()

	return e
}
