package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type SyntaxList struct {
	Base
	children []Node
}

func NewSyntaxList(children []Node) *SyntaxList {
	var width uint32
	for _, child := range children {
		if child != nil {
			width += child.FullWidth()
		}
	}

	return &SyntaxList{
		Base: Base{
			kind:      syntax.SyntaxList,
			fullWidth: width,
		},
		children: children,
	}
}

func (s *SyntaxList) SlotCount() int {
	return len(s.children)
}

func (s *SyntaxList) Slot(index int) Node {
	if index >= 0 && index < len(s.children) {
		return s.children[index]
	}
	return nil
}
