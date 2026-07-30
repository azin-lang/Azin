package object

type Section struct {
	Name        string
	Data        []byte
	Relocations []*Relocation
}
