package elf

import (
	"debug/elf"
	"encoding/binary"
	"os"

	bin "github.com/azin-lang/Azin/internal/binary"
)

const (
	ElfHeaderSize     = 64
	SectionHeaderSize = 64
)

// WriteRelocatableObject generates a standard ELF 64-bit relocatable object file (.o)
// containing the compiled machine code, enabling it to be linked via gcc/ld.
func WriteRelocatableObject(path string, code *bin.Emitter) error {
	e := bin.New(binary.LittleEndian)

	// -------------------------------------------------------------------------
	// 1. ELF Header (64 bytes for ET_REL)
	// -------------------------------------------------------------------------
	e.String("\x7fELF")                  // e_ident[0..3]: Magic number
	e.U8(uint8(elf.ELFCLASS64))          // e_ident[4]: 64-bit object
	e.U8(uint8(elf.ELFDATA2LSB))         // e_ident[5]: Little-endian
	e.U8(uint8(elf.EV_CURRENT))          // e_ident[6]: Current version
	e.U8(uint8(elf.ELFOSABI_NONE))       // e_ident[7]: System V ABI
	e.U8(0)                              // e_ident[8]: ABI version
	e.Zeroes(elf.EI_NIDENT - elf.EI_PAD) // e_ident[9..15]: Padding

	e.U16(uint16(elf.ET_REL))     // e_type: ET_REL (Relocatable file)
	e.U16(uint16(elf.EM_X86_64))  // e_machine: AMD x86-64
	e.U32(uint32(elf.EV_CURRENT)) // e_version: Current version
	e.U64(0)                      // e_entry: 0 for object files
	e.U64(0)                      // e_phoff: No program headers in ET_REL

	// Reserve section header offset to patch after calculating data sizes
	shOffPatch := e.PlaceholderU64()

	e.U32(0)                 // e_flags: Platform-specific flags
	e.U16(ElfHeaderSize)     // e_ehsize: ELF header size
	e.U16(0)                 // e_phentsize: Program header entry size (0)
	e.U16(0)                 // e_phnum: Number of program headers (0)
	e.U16(SectionHeaderSize) // e_shentsize: Section header entry size
	e.U16(5)                 // e_shnum: 5 sections (Null, .text, .symtab, .strtab, .shstrndx)
	e.U16(4)                 // e_shstrndx: Index of section name string table

	// -------------------------------------------------------------------------
	// 2. Section Data Payload
	// -------------------------------------------------------------------------
	// Section 1: .text (Machine Code Payload)
	textOffset := e.Offset()
	textData := code.Data()
	textLen := uint64(len(textData))
	e.Append(code)

	// Section 2: .symtab (Symbol Table Placeholder / Minimal Symbol)
	// Object files need a symbol table so external linkers know what functions are defined.
	symtabOffset := e.Offset()
	// Write a null symbol entry (required by ELF spec)
	e.Zeroes(24)
	// Write symbol entry for 'main'
	// st_name (4), st_info (1), st_other (1), st_shndx (2), st_value (8), st_size (8)
	e.U32(1)       // st_name: offset into .strtab for "main"
	e.U8(0x12)     // st_info: STB_GLOBAL (1 << 4) | STT_FUNC (2)
	e.U8(0)        // st_other: 0
	e.U16(1)       // st_shndx: defined in section 1 (.text)
	e.U64(0)       // st_value: offset 0 in .text
	e.U64(textLen) // st_size
	symtabLen := uint64(e.Offset() - symtabOffset)

	// Section 3: .strtab (String Table for Symbols)
	strtabOffset := e.Offset()
	e.U8(0) // Null string
	e.String("main\x00")
	strtabLen := uint64(e.Offset() - strtabOffset)

	// Section 4: .shstrndx (Section Header String Table)
	shstrtabOffset := e.Offset()
	e.U8(0)
	e.String(".text\x00")
	e.String(".symtab\x00")
	e.String(".strtab\x00")
	e.String(".shstrndx\x00")
	shstrtabLen := uint64(e.Offset() - shstrtabOffset)

	shOffset := e.Offset()
	e.PatchU64(shOffPatch, uint64(shOffset))

	// -------------------------------------------------------------------------
	// 3. Section Header Table (Elf64_Shdr - 5 entries x 64 bytes)
	// -------------------------------------------------------------------------

	// Entry 0: Null Section
	e.U32(0)
	e.U32(0)
	e.U64(0)
	e.U64(0)
	e.U64(0)
	e.U64(0)
	e.U32(0)
	e.U32(0)
	e.U64(0)
	e.U64(0)

	// Entry 1: .text Section
	e.U32(1)                                         // sh_name: offset in .shstrndx ("\n.text") -> offset 1
	e.U32(uint32(elf.SHT_PROGBITS))                  // sh_type: SHT_PROGBITS
	e.U64(uint64(elf.SHF_ALLOC | elf.SHF_EXECINSTR)) // sh_flags: alloc + executable
	e.U64(0)                                         // sh_addr: 0 for relocatable objects
	e.U64(uint64(textOffset))                        // sh_offset
	e.U64(textLen)                                   // sh_size
	e.U32(0)
	e.U32(0)  // sh_link, sh_info
	e.U64(16) // sh_addralign: 16-byte alignment
	e.U64(0)  // sh_entsize

	// Entry 2: .symtab Section
	e.U32(7)                      // sh_name: ".symtab"
	e.U32(uint32(elf.SHT_SYMTAB)) // sh_type
	e.U64(0)                      // sh_flags
	e.U64(0)
	e.U64(uint64(symtabOffset))
	e.U64(symtabLen)
	e.U32(3) // sh_link: index of associated string table (.strtab is index 3)
	e.U32(1) // sh_info: one local/global symbol index threshold
	e.U64(8)
	e.U64(24) // sh_addralign, sh_entsize

	// Entry 3: .strtab Section
	e.U32(15)                     // sh_name: ".strtab"
	e.U32(uint32(elf.SHT_STRTAB)) // sh_type
	e.U64(0)
	e.U64(0)
	e.U64(uint64(strtabOffset))
	e.U64(strtabLen)
	e.U32(0)
	e.U32(0)
	e.U64(1)
	e.U64(0)

	// Entry 4: .shstrndx Section
	e.U32(23)                     // sh_name: ".shstrndx"
	e.U32(uint32(elf.SHT_STRTAB)) // sh_type
	e.U64(0)
	e.U64(0)
	e.U64(uint64(shstrtabOffset))
	e.U64(shstrtabLen)
	e.U32(0)
	e.U32(0)
	e.U64(1)
	e.U64(0)

	// Write object file to disk with standard permissions
	return os.WriteFile(path, e.Data(), 0644)
}
