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
	if greenRoot == nil {
		return nil
	}
	return &Node{
		greenNode: greenRoot,
		parent:    nil,
		position:  0,
	}
}

func (n *Node) Kind() syntax.SyntaxKind {
	if n == nil || n.greenNode == nil {
		return syntax.Unknown
	}
	return n.greenNode.Kind()
}

func (n *Node) Parent() *Node {
	if n == nil {
		return nil
	}
	return n.parent
}

func (n *Node) Position() uint32 {
	if n == nil {
		return 0
	}
	return n.position
}

func (n *Node) FullWidth() uint32 {
	if n == nil || n.greenNode == nil {
		return 0
	}
	return n.greenNode.FullWidth()
}

func (n *Node) Green() green.Node {
	if n == nil {
		return nil
	}
	return n.greenNode
}

// Child lazily evaluates and creates the red wrapper for a green child at a specific index.
func (n *Node) Child(index int) *Node {
	if n == nil || n.greenNode == nil {
		return nil
	}

	greenChild := n.greenNode.Slot(index)
	if greenChild == nil {
		return nil
	}

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
