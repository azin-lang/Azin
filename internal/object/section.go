package object

type SectionKind uint8

const (
	SectionUnknown SectionKind = iota

	// Executable code
	SectionText

	// Read-only initialized data
	SectionRodata

	// Writable initialized data
	SectionData

	// Writable zero-initialized data
	SectionBSS

	// Thread-local initialized data
	SectionTLSData

	// Thread-local zero-initialized data
	SectionTLSBSS

	// Debug information
	SectionDebug

	// Exception handling / unwind tables
	SectionException

	// Symbol/string tables and other metadata
	SectionMetadata
)

type SectionFlags uint32

const (
	SectionReadable SectionFlags = 1 << iota
	SectionWritable
	SectionExecutable
	SectionAlloc
)

type Section struct {
	Name        string
	Kind        SectionKind
	Flags       SectionFlags
	Alignment   uint64
	Data        []byte
	Relocations []Relocation
}
