package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type UnaryExpression struct {
	*Node
}

func AsUnaryExpression(n *Node) *UnaryExpression {
	if n == nil || n.Kind() != syntax.UnaryExpression {
		return nil
	}
	return &UnaryExpression{Node: n}
}

func (u *UnaryExpression) Operator() *Node {
	return u.Child(0)
}

func (u *UnaryExpression) Operand() *Node {
	return u.Child(1)
}
