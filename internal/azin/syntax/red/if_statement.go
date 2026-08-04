package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

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
