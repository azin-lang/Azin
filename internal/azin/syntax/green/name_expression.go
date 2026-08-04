package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type NameExpression struct {
	Base
	identifier *Token
}

func NewNameExpression(identifier *Token) *NameExpression {
	width, flags := ComputeProperties1(identifier)
	return &NameExpression{
		Base: Base{
			kind:      syntax.NameExpression,
			fullWidth: width,
			flags:     flags,
		},
		identifier: identifier,
	}
}

func (n *NameExpression) Identifier() *Token { return n.identifier }
func (n *NameExpression) SlotCount() int     { return 1 }
func (n *NameExpression) Slot(index int) Node {
	if index == 0 && n.identifier != nil {
		return n.identifier
	}
	return nil
}
