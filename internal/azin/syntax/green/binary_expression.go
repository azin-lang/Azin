package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type BinaryExpression struct {
	Base
	left     Node
	operator *Token
	right    Node
}

func NewBinaryExpression(left Node, operator *Token, right Node) *BinaryExpression {
	width, flags := ComputeProperties3(left, operator, right)
	return &BinaryExpression{
		Base: Base{
			kind:      syntax.BinaryExpression,
			fullWidth: width,
			flags:     flags,
		},
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (b *BinaryExpression) Left() Node       { return b.left }
func (b *BinaryExpression) Operator() *Token { return b.operator }
func (b *BinaryExpression) Right() Node      { return b.right }

func (b *BinaryExpression) SlotCount() int { return 3 }
func (b *BinaryExpression) Slot(index int) Node {
	switch index {
	case 0:
		if b.left == nil {
			return nil
		}
		return b.left
	case 1:
		if b.operator == nil {
			return nil
		}
		return b.operator
	case 2:
		if b.right == nil {
			return nil
		}
		return b.right
	default:
		return nil
	}
}
