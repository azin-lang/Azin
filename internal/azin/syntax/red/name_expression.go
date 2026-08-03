package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type NameExpression struct {
	*Node
}

func AsNameExpression(n *Node) *NameExpression {
	if n == nil || n.Kind() != syntax.NameExpression {
		return nil
	}
	return &NameExpression{Node: n}
}

func (n *NameExpression) Identifier() *Node {
	return n.Child(0)
}
