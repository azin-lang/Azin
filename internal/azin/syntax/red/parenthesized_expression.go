package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type ParenthesizedExpression struct {
	*Node
}

func AsParenthesizedExpression(n *Node) *ParenthesizedExpression {
	if n == nil || n.Kind() != syntax.ParenthesizedExpression {
		return nil
	}
	return &ParenthesizedExpression{Node: n}
}

func (p *ParenthesizedExpression) OpenParen() SyntaxToken {
	return p.ChildToken(0)
}

func (p *ParenthesizedExpression) Expression() *Node {
	return p.Child(1)
}

func (p *ParenthesizedExpression) CloseParen() SyntaxToken {
	return p.ChildToken(2)
}
