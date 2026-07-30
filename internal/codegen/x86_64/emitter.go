package x86_64

import (
	binary2 "encoding/binary"

	"github.com/azin-lang/Azin/internal/binary"
)

// EmitMinimalExit returns raw 64-bit x86_64 machine code for:
//
//	mov rax, 60  (sys_exit)
//	mov rdi, 42  (exit code)
//	syscall
func EmitMinimalExit() *binary.Emitter {
	e := binary.New(binary.ByteOrder(binary2.LittleEndian))

	// mov eax, 60
	e.U8(0xB8)
	e.U32(60)

	// mov edi, 42
	e.U8(0xBF)
	e.U32(42)

	// syscall
	e.U8(0x0F)
	e.U8(0x05)

	return e
}
