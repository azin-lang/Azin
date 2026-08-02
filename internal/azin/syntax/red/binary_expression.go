package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

// BinaryExpression provides a typed view over a generic red Node.
type BinaryExpression struct {
	*Node
}

// AsBinaryExpression safely casts a generic Node to a BinaryExpression.
func AsBinaryExpression(n *Node) *BinaryExpression {
	if n == nil || n.Kind() != syntax.BinaryExpression {
		return nil
	}
	return &BinaryExpression{Node: n}
}

func (b *BinaryExpression) Left() *Node {
	return b.Child(0)
}

func (b *BinaryExpression) Operator() *Node {
	return b.Child(1)
}

func (b *BinaryExpression) Right() *Node {
	return b.Child(2)
}
