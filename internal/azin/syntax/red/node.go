package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
)

// Node is the user-facing syntax node. It knows its parent and absolute position.
type Node struct {
	greenNode green.Node
	parent    *Node
	position  uint32 // Absolute position in the source file
}

// NewRoot creates the root of the red tree from a green root node.
func NewRoot(greenRoot green.Node) *Node {
	return &Node{
		greenNode: greenRoot,
		parent:    nil,
		position:  0,
	}
}

func (n *Node) Kind() syntax.SyntaxKind {
	return n.greenNode.Kind()
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) Position() uint32 {
	return n.position
}

func (n *Node) FullWidth() uint32 {
	return n.greenNode.FullWidth()
}

func (n *Node) Green() green.Node {
	return n.greenNode
}

// Child lazily evaluates and creates the red wrapper for a green child at a specific index.
func (n *Node) Child(index int) *Node {
	greenChild := n.greenNode.Slot(index)
	if greenChild == nil {
		return nil
	}

	// Calculate the absolute position of this child by accumulating the widths of all preceding sibling slots.
	childPos := n.position
	for i := range index {
		sibling := n.greenNode.Slot(i)
		if sibling != nil {
			childPos += sibling.FullWidth()
		}
	}

	return &Node{
		greenNode: greenChild,
		parent:    n,
		position:  childPos,
	}
}
