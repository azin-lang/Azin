package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type NameExpression struct {
	Base
	identifier *Token
}

func NewNameExpression(identifier *Token) *NameExpression {
	var width uint32
	if identifier != nil {
		width = identifier.FullWidth()
	}

	return &NameExpression{
		Base: Base{
			kind:      syntax.NameExpression,
			fullWidth: width,
		},
		identifier: identifier,
	}
}

func (n *NameExpression) Identifier() *Token { return n.identifier }
func (n *NameExpression) SlotCount() int     { return 1 }
func (n *NameExpression) Slot(index int) Node {
	if index == 0 {
		return n.identifier
	}
	return nil
}
