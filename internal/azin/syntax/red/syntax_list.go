package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type SyntaxList struct {
	*Node
}

func AsSyntaxList(n *Node) *SyntaxList {
	if n == nil || n.Kind() != syntax.SyntaxList {
		return nil
	}
	return &SyntaxList{Node: n}
}

func (s *SyntaxList) Count() int {
	if s == nil || s.greenNode == nil {
		return 0
	}
	return s.greenNode.SlotCount()
}

func (s *SyntaxList) Item(index int) *Node {
	if index < 0 || index >= s.Count() {
		return nil
	}
	return s.Child(index)
}
