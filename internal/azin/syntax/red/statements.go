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

type IfStatement struct {
	*Node
}

func AsIfStatement(n *Node) *IfStatement {
	if n == nil || n.Kind() != syntax.IfStatement {
		return nil
	}
	return &IfStatement{Node: n}
}

func (i *IfStatement) IfKeyword() SyntaxToken   { return i.ChildToken(0) }
func (i *IfStatement) Condition() *Node         { return i.Child(1) }
func (i *IfStatement) ThenBranch() *Node        { return i.Child(2) }
func (i *IfStatement) ElseKeyword() SyntaxToken { return i.ChildToken(3) }
func (i *IfStatement) ElseBranch() *Node        { return i.Child(4) }

type ReturnStatement struct {
	*Node
}

func AsReturnStatement(n *Node) *ReturnStatement {
	if n == nil || n.Kind() != syntax.ReturnStatement {
		return nil
	}
	return &ReturnStatement{Node: n}
}

func (r *ReturnStatement) ReturnKeyword() SyntaxToken { return r.ChildToken(0) }
func (r *ReturnStatement) Expression() *Node          { return r.Child(1) }
