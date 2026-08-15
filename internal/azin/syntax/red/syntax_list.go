package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
)

type SyntaxList struct {
	*Node
}

func AsSyntaxList(n *Node) *SyntaxList {
	if n == nil || n.Kind() != syntax.SyntaxList {
		return nil
	}
	return &SyntaxList{Node: n}
}

// Count returns the total number of items, transparently flattening balanced green trees.
func (s *SyntaxList) Count() int {
	if s == nil || green.IsNil(s.greenNode) {
		return 0
	}
	return countGreenListItems(s.greenNode)
}

func countGreenListItems(g green.Node) int {
	if green.IsNil(g) {
		return 0
	}
	if g.Kind() != syntax.SyntaxList {
		return 1
	}
	total := 0
	for i := range g.SlotCount() {
		total += countGreenListItems(g.Slot(i))
	}
	return total
}

// Item retrieves the i-th logical element node in the list.
func (s *SyntaxList) Item(index int) *Node {
	if s == nil || index < 0 || index >= s.Count() {
		return nil
	}
	return getItemFromList(s.Node, index)
}

func getItemFromList(n *Node, index int) *Node {
	if n == nil || green.IsNil(n.greenNode) {
		return nil
	}

	currentIndex := index
	for i := range n.greenNode.SlotCount() {
		slot := n.greenNode.Slot(i)
		if green.IsNil(slot) {
			continue
		}

		itemCount := countGreenListItems(slot)
		if currentIndex < itemCount {
			childNode := n.Child(i)
			if slot.Kind() == syntax.SyntaxList {
				return getItemFromList(childNode, currentIndex)
			}
			return childNode
		}
		currentIndex -= itemCount
	}

	return nil
}
