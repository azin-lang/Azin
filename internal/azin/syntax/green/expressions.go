package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type BinaryExpression struct {
	Base
	children [3]Node
}

func NewBinaryExpression(left Node, operator *Token, right Node) *BinaryExpression {
	children := [3]Node{NilSafe(left), NilSafe(operator), NilSafe(right)}
	width, flags := ComputeProperties(children[:])
	return &BinaryExpression{
		Base:     Base{kind: syntax.BinaryExpression, fullWidth: width, flags: flags},
		children: children,
	}
}

func (b *BinaryExpression) Left() Node          { return b.children[0] }
func (b *BinaryExpression) Operator() *Token    { t, _ := b.children[1].(*Token); return t }
func (b *BinaryExpression) Right() Node         { return b.children[2] }
func (b *BinaryExpression) IsNil() bool         { return b == nil }
func (b *BinaryExpression) SlotCount() int      { return 3 }
func (b *BinaryExpression) Slot(index int) Node { return b.children[index] }

type CallExpression struct {
	Base
	children [4]Node
}

func NewCallExpression(expression Node, openParen *Token, arguments Node, closeParen *Token) *CallExpression {
	children := [4]Node{NilSafe(expression), NilSafe(openParen), NilSafe(arguments), NilSafe(closeParen)}
	width, flags := ComputeProperties(children[:])
	return &CallExpression{
		Base:     Base{kind: syntax.CallExpression, fullWidth: width, flags: flags},
		children: children,
	}
}

func (c *CallExpression) Expression() Node    { return c.children[0] }
func (c *CallExpression) OpenParen() *Token   { t, _ := c.children[1].(*Token); return t }
func (c *CallExpression) Arguments() Node     { return c.children[2] }
func (c *CallExpression) CloseParen() *Token  { t, _ := c.children[3].(*Token); return t }
func (c *CallExpression) IsNil() bool         { return c == nil }
func (c *CallExpression) SlotCount() int      { return 4 }
func (c *CallExpression) Slot(index int) Node { return c.children[index] }

type LiteralExpression struct {
	Base
	children [1]Node
}

func NewLiteralExpression(token *Token) *LiteralExpression {
	children := [1]Node{NilSafe(token)}
	width, flags := ComputeProperties(children[:])
	return &LiteralExpression{
		Base:     Base{kind: syntax.LiteralExpression, fullWidth: width, flags: flags},
		children: children,
	}
}

func (l *LiteralExpression) Token() *Token       { t, _ := l.children[0].(*Token); return t }
func (l *LiteralExpression) IsNil() bool         { return l == nil }
func (l *LiteralExpression) SlotCount() int      { return 1 }
func (l *LiteralExpression) Slot(index int) Node { return l.children[index] }

type MemberAccessExpression struct {
	Base
	children [3]Node
}

func NewMemberAccessExpression(expression Node, dot, name *Token) *MemberAccessExpression {
	children := [3]Node{NilSafe(expression), NilSafe(dot), NilSafe(name)}
	width, flags := ComputeProperties(children[:])
	return &MemberAccessExpression{
		Base:     Base{kind: syntax.MemberAccessExpression, fullWidth: width, flags: flags},
		children: children,
	}
}

func (m *MemberAccessExpression) Expression() Node    { return m.children[0] }
func (m *MemberAccessExpression) Dot() *Token         { t, _ := m.children[1].(*Token); return t }
func (m *MemberAccessExpression) Name() *Token        { t, _ := m.children[2].(*Token); return t }
func (m *MemberAccessExpression) IsNil() bool         { return m == nil }
func (m *MemberAccessExpression) SlotCount() int      { return 3 }
func (m *MemberAccessExpression) Slot(index int) Node { return m.children[index] }

type NameExpression struct {
	Base
	children [1]Node
}

func NewNameExpression(identifier *Token) *NameExpression {
	children := [1]Node{NilSafe(identifier)}
	width, flags := ComputeProperties(children[:])
	return &NameExpression{
		Base:     Base{kind: syntax.NameExpression, fullWidth: width, flags: flags},
		children: children,
	}
}

func (n *NameExpression) Identifier() *Token  { t, _ := n.children[0].(*Token); return t }
func (n *NameExpression) IsNil() bool         { return n == nil }
func (n *NameExpression) SlotCount() int      { return 1 }
func (n *NameExpression) Slot(index int) Node { return n.children[index] }

type ParenthesizedExpression struct {
	Base
	children [3]Node
}

func NewParenthesizedExpression(openParen *Token, expression Node, closeParen *Token) *ParenthesizedExpression {
	children := [3]Node{NilSafe(openParen), NilSafe(expression), NilSafe(closeParen)}
	width, flags := ComputeProperties(children[:])
	return &ParenthesizedExpression{
		Base:     Base{kind: syntax.ParenthesizedExpression, fullWidth: width, flags: flags},
		children: children,
	}
}

func (p *ParenthesizedExpression) OpenParen() *Token   { t, _ := p.children[0].(*Token); return t }
func (p *ParenthesizedExpression) Expression() Node    { return p.children[1] }
func (p *ParenthesizedExpression) CloseParen() *Token  { t, _ := p.children[2].(*Token); return t }
func (p *ParenthesizedExpression) IsNil() bool         { return p == nil }
func (p *ParenthesizedExpression) SlotCount() int      { return 3 }
func (p *ParenthesizedExpression) Slot(index int) Node { return p.children[index] }

type UnaryExpression struct {
	Base
	children [2]Node
}

func NewUnaryExpression(operator *Token, operand Node) *UnaryExpression {
	children := [2]Node{NilSafe(operator), NilSafe(operand)}
	width, flags := ComputeProperties(children[:])
	return &UnaryExpression{
		Base:     Base{kind: syntax.UnaryExpression, fullWidth: width, flags: flags},
		children: children,
	}
}

func (u *UnaryExpression) Operator() *Token    { t, _ := u.children[0].(*Token); return t }
func (u *UnaryExpression) Operand() Node       { return u.children[1] }
func (u *UnaryExpression) IsNil() bool         { return u == nil }
func (u *UnaryExpression) SlotCount() int      { return 2 }
func (u *UnaryExpression) Slot(index int) Node { return u.children[index] }
