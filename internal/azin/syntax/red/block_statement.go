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

func (b *BlockStatement) OpenBrace() SyntaxToken  { return b.ChildToken(0) }
func (b *BlockStatement) Statements() *SyntaxList { return AsSyntaxList(b.Child(1)) }
func (b *BlockStatement) CloseBrace() SyntaxToken { return b.ChildToken(2) }
