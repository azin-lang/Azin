package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type CallExpression struct {
	Base
	expression Node
	openParen  *Token
	arguments  Node
	closeParen *Token
}

func NewCallExpression(expression Node, openParen *Token, arguments Node, closeParen *Token) *CallExpression {
	width, flags := ComputeProperties4(expression, openParen, arguments, closeParen)
	return &CallExpression{
		Base: Base{
			kind:      syntax.CallExpression,
			fullWidth: width,
			flags:     flags,
		},
		expression: expression,
		openParen:  openParen,
		arguments:  arguments,
		closeParen: closeParen,
	}
}

func (c *CallExpression) Expression() Node   { return c.expression }
func (c *CallExpression) OpenParen() *Token  { return c.openParen }
func (c *CallExpression) Arguments() Node    { return c.arguments }
func (c *CallExpression) CloseParen() *Token { return c.closeParen }

func (c *CallExpression) SlotCount() int { return 4 }
func (c *CallExpression) Slot(index int) Node {
	switch index {
	case 0:
		if c.expression == nil {
			return nil
		}
		return c.expression
	case 1:
		if c.openParen == nil {
			return nil
		}
		return c.openParen
	case 2:
		if c.arguments == nil {
			return nil
		}
		return c.arguments
	case 3:
		if c.closeParen == nil {
			return nil
		}
		return c.closeParen
	default:
		return nil
	}
}
