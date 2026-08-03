package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type UnaryExpression struct {
	Base

	operator *Token
	operand  Node
}

func NewUnaryExpression(operator *Token, operand Node) *UnaryExpression {
	var width uint32

	if operator != nil {
		width += operator.FullWidth()
	}

	if operand != nil {
		width += operand.FullWidth()
	}

	return &UnaryExpression{
		Base: Base{
			kind:      syntax.UnaryExpression,
			fullWidth: width,
		},
		operator: operator,
		operand:  operand,
	}
}

func (u *UnaryExpression) Operator() *Token { return u.operator }
func (u *UnaryExpression) Operand() Node    { return u.operand }

func (u *UnaryExpression) SlotCount() int { return 2 }

func (u *UnaryExpression) Slot(index int) Node {
	switch index {
	case 0:
		return u.operator
	case 1:
		return u.operand
	default:
		return nil
	}
}
