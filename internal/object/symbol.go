package object

type Symbol struct {
	Name    string
	Section *Section
	Offset  uint64
}
