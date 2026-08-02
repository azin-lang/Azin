package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

// BinaryExpression represents a node with a left side, an operator, and a right side.
type BinaryExpression struct {
	Base

	left     Node
	operator *Token
	right    Node
}

func NewBinaryExpression(left Node, operator *Token, right Node) *BinaryExpression {
	var width uint32
	if left != nil {
		width += left.FullWidth()
	}
	if operator != nil {
		width += operator.FullWidth()
	}
	if right != nil {
		width += right.FullWidth()
	}

	return &BinaryExpression{
		Base: Base{
			kind:      syntax.BinaryExpression,
			fullWidth: width,
		},
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (b *BinaryExpression) Left() Node {
	return b.left
}

func (b *BinaryExpression) Operator() *Token {
	return b.operator
}

func (b *BinaryExpression) Right() Node {
	return b.right
}

// SlotCount dictates exactly how many children this node expects.
func (b *BinaryExpression) SlotCount() int {
	return 3
}

// Slot returns the child at the given index. This allows generic tree traversal.
func (b *BinaryExpression) Slot(index int) Node {
	switch index {
	case 0:
		return b.left
	case 1:
		return b.operator
	case 2:
		return b.right
	default:
		return nil
	}
}
