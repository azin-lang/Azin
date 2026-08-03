package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type BlockStatement struct {
	*Node
}

func AsBlockStatement(n *Node) *BlockStatement {
	if n == nil || n.Kind() != syntax.BlockStatement {
		return nil
	}
	return &BlockStatement{Node: n}
}

func (b *BlockStatement) OpenBrace() *Node  { return b.Child(0) }
func (b *BlockStatement) Statements() *Node { return b.Child(1) }
func (b *BlockStatement) CloseBrace() *Node { return b.Child(2) }
