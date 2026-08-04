package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type MemberAccessExpression struct {
	*Node
}

func AsMemberAccessExpression(n *Node) *MemberAccessExpression {
	if n == nil || n.Kind() != syntax.MemberAccessExpression {
		return nil
	}
	return &MemberAccessExpression{Node: n}
}

func (m *MemberAccessExpression) Expression() *Node { return m.Child(0) }
func (m *MemberAccessExpression) Dot() SyntaxToken  { return m.ChildToken(1) }
func (m *MemberAccessExpression) Name() SyntaxToken { return m.ChildToken(2) }
