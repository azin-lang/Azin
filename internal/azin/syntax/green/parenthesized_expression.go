package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type ParenthesizedExpression struct {
	Base
	openParen  *Token
	expression Node
	closeParen *Token
}

func NewParenthesizedExpression(openParen *Token, expression Node, closeParen *Token) *ParenthesizedExpression {
	width, flags := ComputeProperties3(openParen, expression, closeParen)
	return &ParenthesizedExpression{
		Base: Base{
			kind:      syntax.ParenthesizedExpression,
			fullWidth: width,
			flags:     flags,
		},
		openParen:  openParen,
		expression: expression,
		closeParen: closeParen,
	}
}

func (p *ParenthesizedExpression) OpenParen() *Token  { return p.openParen }
func (p *ParenthesizedExpression) Expression() Node   { return p.expression }
func (p *ParenthesizedExpression) CloseParen() *Token { return p.closeParen }

func (p *ParenthesizedExpression) SlotCount() int { return 3 }
func (p *ParenthesizedExpression) Slot(index int) Node {
	switch index {
	case 0:
		if p.openParen == nil {
			return nil
		}
		return p.openParen
	case 1:
		if p.expression == nil {
			return nil
		}
		return p.expression
	case 2:
		if p.closeParen == nil {
			return nil
		}
		return p.closeParen
	default:
		return nil
	}
}
