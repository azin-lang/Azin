package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type MemberAccessExpression struct {
	Base
	expression Node
	dot        *Token
	name       *Token
}

func NewMemberAccessExpression(expression Node, dot *Token, name *Token) *MemberAccessExpression {
	width, flags := ComputeProperties3(expression, dot, name)
	return &MemberAccessExpression{
		Base: Base{
			kind:      syntax.MemberAccessExpression,
			fullWidth: width,
			flags:     flags,
		},
		expression: expression,
		dot:        dot,
		name:       name,
	}
}

func (m *MemberAccessExpression) Expression() Node { return m.expression }
func (m *MemberAccessExpression) Dot() *Token      { return m.dot }
func (m *MemberAccessExpression) Name() *Token     { return m.name }

func (m *MemberAccessExpression) SlotCount() int { return 3 }
func (m *MemberAccessExpression) Slot(index int) Node {
	switch index {
	case 0:
		if m.expression == nil {
			return nil
		}
		return m.expression
	case 1:
		if m.dot == nil {
			return nil
		}
		return m.dot
	case 2:
		if m.name == nil {
			return nil
		}
		return m.name
	default:
		return nil
	}
}
