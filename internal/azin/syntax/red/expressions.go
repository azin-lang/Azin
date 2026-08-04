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

type CallExpression struct {
	*Node
}

func AsCallExpression(n *Node) *CallExpression {
	if n == nil || n.Kind() != syntax.CallExpression {
		return nil
	}
	return &CallExpression{Node: n}
}

func (c *CallExpression) Expression() *Node       { return c.Child(0) }
func (c *CallExpression) OpenParen() SyntaxToken  { return c.ChildToken(1) }
func (c *CallExpression) Arguments() *SyntaxList  { return AsSyntaxList(c.Child(2)) }
func (c *CallExpression) CloseParen() SyntaxToken { return c.ChildToken(3) }

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

type NameExpression struct {
	*Node
}

func AsNameExpression(n *Node) *NameExpression {
	if n == nil || n.Kind() != syntax.NameExpression {
		return nil
	}
	return &NameExpression{Node: n}
}

func (n *NameExpression) Identifier() SyntaxToken {
	return n.ChildToken(0)
}

type ParenthesizedExpression struct {
	*Node
}

func AsParenthesizedExpression(n *Node) *ParenthesizedExpression {
	if n == nil || n.Kind() != syntax.ParenthesizedExpression {
		return nil
	}
	return &ParenthesizedExpression{Node: n}
}

func (p *ParenthesizedExpression) OpenParen() SyntaxToken {
	return p.ChildToken(0)
}

func (p *ParenthesizedExpression) Expression() *Node {
	return p.Child(1)
}

func (p *ParenthesizedExpression) CloseParen() SyntaxToken {
	return p.ChildToken(2)
}

type UnaryExpression struct {
	*Node
}

func AsUnaryExpression(n *Node) *UnaryExpression {
	if n == nil || n.Kind() != syntax.UnaryExpression {
		return nil
	}
	return &UnaryExpression{Node: n}
}

func (u *UnaryExpression) Operator() SyntaxToken {
	return u.ChildToken(0)
}

func (u *UnaryExpression) Operand() *Node {
	return u.Child(1)
}
