package object

type SymbolKind uint8

const (
	SymbolUnknown SymbolKind = iota

	SymbolFunction
	SymbolObject
	SymbolSection
	SymbolFile
)

type SymbolBinding uint8

const (
	BindingLocal SymbolBinding = iota
	BindingGlobal
	BindingWeak
)

type Symbol struct {
	Name    string
	Kind    SymbolKind
	Binding SymbolBinding

	Section *Section
	Offset  uint64
	Size    uint64
}
