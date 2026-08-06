package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/text"
)

// Node is the user-facing syntax node. It knows its parent and absolute position.
type Node struct {
	greenNode green.Node
	parent    *Node
	position  uint32
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

func (n *Node) Kind() syntax.Kind {
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

func (n *Node) FullSpan() text.Span {
	if n == nil || green.IsNil(n.greenNode) {
		return text.Span{}
	}
	return text.NewSpan(n.position, n.greenNode.FullWidth())
}

func (n *Node) Span() text.Span {
	if n == nil || green.IsNil(n.greenNode) {
		return text.Span{}
	}

	leading := green.GetLeadingTriviaWidth(n.greenNode)
	trailing := green.GetTrailingTriviaWidth(n.greenNode)

	start := n.position + leading
	length := n.greenNode.FullWidth() - leading - trailing
	return text.NewSpan(start, length)
}

func (n *Node) ChildToken(index int) SyntaxToken {
	if n == nil || green.IsNil(n.greenNode) || index < 0 || index >= n.greenNode.SlotCount() {
		return SyntaxToken{}
	}

	greenChild := n.greenNode.Slot(index)
	if green.IsNil(greenChild) {
		return SyntaxToken{}
	}

	greenTok, ok := greenChild.(*green.Token)
	if !ok {
		return SyntaxToken{}
	}

	childPos := n.position
	for i := range index {
		if sibling := n.greenNode.Slot(i); !green.IsNil(sibling) {
			childPos += sibling.FullWidth()
		}
	}

	return NewSyntaxToken(n, greenTok, childPos, index)
}
