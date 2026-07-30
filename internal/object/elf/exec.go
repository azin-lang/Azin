package elf

import (
	"debug/elf"
	"encoding/binary"
	"os"

	bin "github.com/azin-lang/Azin/internal/binary"
)

const (
	ElfHeaderSize  = 64
	ProgHeaderSize = 56
	LoadAddr       = 0x400000 // Standard 64-bit ELF base address
)

func WriteMinimalExecutable(path string, code *bin.Emitter) error {
	e := bin.New(binary.LittleEndian)

	e.String("\x7fELF")                  // e_ident[0..3]: Magic number
	e.U8(uint8(elf.ELFCLASS64))          // e_ident[4] (EI_CLASS): 64-bit
	e.U8(uint8(elf.ELFDATA2LSB))         // e_ident[5] (EI_DATA): Little Endian
	e.U8(uint8(elf.EV_CURRENT))          // e_ident[6] (EI_VERSION): Current version
	e.U8(uint8(elf.ELFOSABI_LINUX))      // e_ident[7] (EI_OSABI): Linux
	e.U8(0)                              // e_ident[8] (EI_ABIVERSION): 0
	e.Zeroes(elf.EI_NIDENT - elf.EI_PAD) // e_ident[9..15] (EI_PAD): 7 padding bytes (16 - 9)

	e.U16(uint16(elf.ET_EXEC))    // e_type: ET_EXEC
	e.U16(uint16(elf.EM_X86_64))  // e_machine: EM_X86_64
	e.U32(uint32(elf.EV_CURRENT)) // e_version: EV_CURRENT

	// Reserve e_entry offset to patch after calculating entry point address
	entryPatchOff := e.PlaceholderU64()

	e.U64(ElfHeaderSize)  // e_phoff: Program header table offset immediately follows ELF Header
	e.U64(0)              // e_shoff: No section headers in this minimal binary
	e.U32(0)              // e_flags
	e.U16(ElfHeaderSize)  // e_ehsize: Size of ELF header
	e.U16(ProgHeaderSize) // e_phentsize: Size of program header entry
	e.U16(1)              // e_phnum: 1 Program header
	e.U16(0)              // e_shentsize
	e.U16(0)              // e_shnum
	e.U16(0)              // e_shstrndx

	e.U32(uint32(elf.PT_LOAD))         // p_type: PT_LOAD
	e.U32(uint32(elf.PF_R | elf.PF_X)) // p_flags: PF_R | PF_X (Readable + Executable)
	e.U64(0)                           // p_offset: Segment starts at file start
	e.U64(LoadAddr)                    // p_vaddr: Virtual load address
	e.U64(LoadAddr)                    // p_paddr: Physical load address

	// Reserve size offsets to patch once all code is appended
	fileSizePatchOff := e.PlaceholderU64()
	memSizePatchOff := e.PlaceholderU64()

	e.U64(0x1000) // p_align: 4KB page alignment

	entryAddr := LoadAddr + e.Offset()
	e.Append(code)

	totalSize := e.Size()

	// Patch entry point and segment sizes back into the header
	e.PatchU64(entryPatchOff, entryAddr)
	e.PatchU64(fileSizePatchOff, totalSize)
	e.PatchU64(memSizePatchOff, totalSize)

	// Write executable to disk with 0755 (+x) permissions
	return os.WriteFile(path, e.Data(), 0755)
}
