package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type BinaryExpression struct {
	*Node
}

func AsBinaryExpression(n *Node) *BinaryExpression {
	if n == nil || n.Kind() != syntax.BinaryExpression {
		return nil
	}
	return &BinaryExpression{Node: n}
}

func (b *BinaryExpression) Left() *Node           { return b.Child(0) }
func (b *BinaryExpression) Operator() SyntaxToken { return b.ChildToken(1) }
func (b *BinaryExpression) Right() *Node          { return b.Child(2) }
