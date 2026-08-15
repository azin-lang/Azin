package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/text"
)

// Node is the user-facing red-tree node.
//
// Green nodes contain immutable structural information. Red nodes add
// contextual information such as the parent and absolute source position.
//
// Red nodes are created lazily during traversal.
type Node struct {
	greenNode green.Node
	parent    *Node
	position  int
}

// NewRoot creates a red node around a green root node.
func NewRoot(greenRoot green.Node) *Node {
	if green.IsNil(greenRoot) {
		return nil
	}

	return &Node{
		greenNode: greenRoot,
		position:  0,
	}
}

// Green returns the underlying immutable green node.
func (n *Node) Green() green.Node {
	if n == nil {
		return nil
	}

	return n.greenNode
}

// Parent returns the containing red node.
//
// The root node has no parent.
func (n *Node) Parent() *Node {
	if n == nil {
		return nil
	}

	return n.parent
}

// Kind returns the syntax kind of the node.
func (n *Node) Kind() syntax.Kind {
	if n == nil || green.IsNil(n.greenNode) {
		return syntax.Unknown
	}

	return n.greenNode.Kind()
}

// Position returns the absolute byte offset of the beginning of the node's
// full span.
func (n *Node) Position() int {
	if n == nil {
		return 0
	}

	return n.position
}

// FullWidth returns the width of the node including trivia.
func (n *Node) FullWidth() int {
	if n == nil || green.IsNil(n.greenNode) {
		return 0
	}

	return int(n.greenNode.FullWidth())
}

// SlotCount returns the number of slots in the underlying green node.
func (n *Node) SlotCount() int {
	if n == nil || green.IsNil(n.greenNode) {
		return 0
	}

	return n.greenNode.SlotCount()
}

// childPosition returns the absolute position of a child slot.
//
// Green nodes store widths rather than absolute positions. The red layer
// reconstructs the absolute position from the parent's position and the
// widths of preceding siblings.
func (n *Node) childPosition(index int) int {
	position := n.position

	for i := range index {
		child := n.greenNode.Slot(i)
		if !green.IsNil(child) {
			position += int(child.FullWidth())
		}
	}

	return position
}

// Child returns a lazily-created red node for the specified slot.
func (n *Node) Child(index int) *Node {
	if n == nil || green.IsNil(n.greenNode) {
		return nil
	}

	if index < 0 || index >= n.greenNode.SlotCount() {
		return nil
	}

	child := n.greenNode.Slot(index)
	if green.IsNil(child) {
		return nil
	}

	return &Node{
		greenNode: child,
		parent:    n,
		position:  n.childPosition(index),
	}
}

// ChildToken returns a red syntax token for the specified slot.
//
// A zero SyntaxToken is returned if the slot does not contain a token.
func (n *Node) ChildToken(index int) SyntaxToken {
	if n == nil || green.IsNil(n.greenNode) {
		return SyntaxToken{}
	}

	if index < 0 || index >= n.greenNode.SlotCount() {
		return SyntaxToken{}
	}

	child := n.greenNode.Slot(index)
	if green.IsNil(child) {
		return SyntaxToken{}
	}

	token, ok := child.(*green.Token)
	if !ok {
		return SyntaxToken{}
	}

	return SyntaxToken{
		parent:   n,
		green:    token,
		position: n.childPosition(index),
		index:    index,
	}
}

// FullSpan returns the complete source span of the node, including trivia.
func (n *Node) FullSpan() text.Span {
	if n == nil || green.IsNil(n.greenNode) {
		return text.Span{}
	}

	return text.NewSpan(
		n.position,
		int(n.greenNode.FullWidth()),
	)
}

// Span returns the source span of the node excluding its leading and trailing
// trivia.
func (n *Node) Span() text.Span {
	if n == nil || green.IsNil(n.greenNode) {
		return text.Span{}
	}

	leading := int(green.GetLeadingTriviaWidth(n.greenNode))
	trailing := int(green.GetTrailingTriviaWidth(n.greenNode))
	fullWidth := int(n.greenNode.FullWidth())

	start := n.position + leading
	width := fullWidth - leading - trailing
	if width < 0 {
		width = 0
	}

	return text.NewSpan(start, width)
}
