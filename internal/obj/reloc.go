package obj

type Relocation struct {
	Offset uint32 // Where the placeholder is sitting in the .text machine code
	Symbol string // What function it wants to call ("printf")
	Type   uint8  // How to patch it (e.g. relative 32-bit jump)
}
