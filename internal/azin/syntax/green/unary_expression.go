package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type UnaryExpression struct {
	Base
	operator *Token
	operand  Node
}

func NewUnaryExpression(operator *Token, operand Node) *UnaryExpression {
	width, flags := ComputeProperties2(operator, operand)
	return &UnaryExpression{
		Base: Base{
			kind:      syntax.UnaryExpression,
			fullWidth: width,
			flags:     flags,
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
		if u.operator == nil {
			return nil
		}
		return u.operator
	case 1:
		if u.operand == nil {
			return nil
		}
		return u.operand
	default:
		return nil
	}
}
