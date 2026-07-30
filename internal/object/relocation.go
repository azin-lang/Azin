package object

type RelocationType uint8

const (
	RelocAbsolute RelocationType = iota
	RelocPCRelative
	RelocCall
	RelocJump
	RelocGOT
	RelotPLT
	RelocTLS
)

type Relocation struct {
	Offset uint64
	Symbol *Symbol
	Type   RelocationType
	Addend int64
}
