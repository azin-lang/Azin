package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type CallExpression struct {
	*Node
}

func AsCallExpression(n *Node) *CallExpression {
	if n == nil || n.Kind() != syntax.CallExpression {
		return nil
	}
	return &CallExpression{Node: n}
}

func (c *CallExpression) Expression() *Node       { return c.Child(0) }
func (c *CallExpression) OpenParen() SyntaxToken  { return c.ChildToken(1) }
func (c *CallExpression) Arguments() *SyntaxList  { return AsSyntaxList(c.Child(2)) }
func (c *CallExpression) CloseParen() SyntaxToken { return c.ChildToken(3) }
