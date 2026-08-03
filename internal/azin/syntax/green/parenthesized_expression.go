package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type ParenthesizedExpression struct {
	Base

	openParen  *Token
	expression Node
	closeParen *Token
}

func NewParenthesizedExpression(openParen *Token, expression Node, closeParen *Token) *ParenthesizedExpression {
	var width uint32
	if openParen != nil {
		width += openParen.FullWidth()
	}
	if expression != nil {
		width += expression.FullWidth()
	}
	if closeParen != nil {
		width += closeParen.FullWidth()
	}

	return &ParenthesizedExpression{
		Base: Base{
			kind:      syntax.ParenthesizedExpression,
			fullWidth: width,
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
		return p.openParen
	case 1:
		return p.expression
	case 2:
		return p.closeParen
	default:
		return nil
	}
}
