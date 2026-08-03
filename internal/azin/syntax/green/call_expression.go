package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type CallExpression struct {
	Base
	expression Node
	openParen  *Token
	arguments  []Node // Comma-separated arguments or SyntaxList
	closeParen *Token
}

func NewCallExpression(expression Node, openParen *Token, arguments []Node, closeParen *Token) *CallExpression {
	var width uint32
	if expression != nil {
		width += expression.FullWidth()
	}
	if openParen != nil {
		width += openParen.FullWidth()
	}
	for _, arg := range arguments {
		if arg != nil {
			width += arg.FullWidth()
		}
	}
	if closeParen != nil {
		width += closeParen.FullWidth()
	}

	return &CallExpression{
		Base: Base{
			kind:      syntax.CallExpression,
			fullWidth: width,
		},
		expression: expression,
		openParen:  openParen,
		arguments:  arguments,
		closeParen: closeParen,
	}
}

func (c *CallExpression) Expression() Node   { return c.expression }
func (c *CallExpression) OpenParen() *Token  { return c.openParen }
func (c *CallExpression) CloseParen() *Token { return c.closeParen }

func (c *CallExpression) SlotCount() int {
	return 3 + len(c.arguments)
}

func (c *CallExpression) Slot(index int) Node {
	switch index {
	case 0:
		return c.expression
	case 1:
		return c.openParen
	default:
		argIndex := index - 2
		if argIndex < len(c.arguments) {
			return c.arguments[argIndex]
		}
		if index == 2+len(c.arguments) {
			return c.closeParen
		}
		return nil
	}
}
