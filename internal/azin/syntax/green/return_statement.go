package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type ReturnStatement struct {
	Base
	returnKeyword *Token
	expression    Node
}

func NewReturnStatement(returnKeyword *Token, expression Node) *ReturnStatement {
	var width uint32
	if returnKeyword != nil {
		width += returnKeyword.FullWidth()
	}
	if expression != nil {
		width += expression.FullWidth()
	}

	return &ReturnStatement{
		Base: Base{
			kind:      syntax.ReturnStatement,
			fullWidth: width,
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
		return r.returnKeyword
	case 1:
		return r.expression
	default:
		return nil
	}
}
