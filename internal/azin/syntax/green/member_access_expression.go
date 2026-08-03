package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type MemberAccessExpression struct {
	Base
	expression Node
	dot        *Token
	name       *Token
}

func NewMemberAccessExpression(expression Node, dot *Token, name *Token) *MemberAccessExpression {
	var width uint32
	if expression != nil {
		width += expression.FullWidth()
	}
	if dot != nil {
		width += dot.FullWidth()
	}
	if name != nil {
		width += name.FullWidth()
	}

	return &MemberAccessExpression{
		Base: Base{
			kind:      syntax.MemberAccessExpression,
			fullWidth: width,
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
		return m.expression
	case 1:
		return m.dot
	case 2:
		return m.name
	default:
		return nil
	}
}
