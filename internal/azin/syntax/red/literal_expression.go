package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type LiteralExpression struct {
	*Node
}

func AsLiteralExpression(n *Node) *LiteralExpression {
	if n == nil || n.Kind() != syntax.LiteralExpression {
		return nil
	}
	return &LiteralExpression{Node: n}
}

func (l *LiteralExpression) Token() SyntaxToken {
	return l.ChildToken(0)
}
