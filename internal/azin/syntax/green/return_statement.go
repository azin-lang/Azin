package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type ReturnStatement struct {
	Base
	returnKeyword *Token
	expression    Node
}

func NewReturnStatement(returnKeyword *Token, expression Node) *ReturnStatement {
	width, flags := ComputeProperties2(returnKeyword, expression)
	return &ReturnStatement{
		Base: Base{
			kind:      syntax.ReturnStatement,
			fullWidth: width,
			flags:     flags,
		},
		returnKeyword: returnKeyword,
		expression:    expression,
	}
}

func (r *ReturnStatement) ReturnKeyword() *Token { return r.returnKeyword }
func (r *ReturnStatement) Expression() Node      { return r.expression }

func (r *ReturnStatement) SlotCount() int { return 2 }
func (r *ReturnStatement) Slot(index int) Node {
	switch index {
	case 0:
		if r.returnKeyword == nil {
			return nil
		}
		return r.returnKeyword
	case 1:
		if r.expression == nil {
			return nil
		}
		return r.expression
	default:
		return nil
	}
}
