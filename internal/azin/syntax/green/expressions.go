package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type BinaryExpression struct {
	Base
	left     Node
	operator *Token
	right    Node
}

func NewBinaryExpression(left Node, operator *Token, right Node) *BinaryExpression {
	width, flags := ComputeProperties3(left, operator, right)
	return &BinaryExpression{
		Base: Base{
			kind:      syntax.BinaryExpression,
			fullWidth: width,
			flags:     flags,
		},
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (b *BinaryExpression) Left() Node       { return b.left }
func (b *BinaryExpression) Operator() *Token { return b.operator }
func (b *BinaryExpression) Right() Node      { return b.right }

func (b *BinaryExpression) IsNil() bool { return b == nil }

func (b *BinaryExpression) SlotCount() int { return 3 }
func (b *BinaryExpression) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(b.left)
	case 1:
		return SafeNode(b.operator)
	case 2:
		return SafeNode(b.right)
	default:
		return nil
	}
}

type CallExpression struct {
	Base
	expression Node
	openParen  *Token
	arguments  Node
	closeParen *Token
}

func NewCallExpression(expression Node, openParen *Token, arguments Node, closeParen *Token) *CallExpression {
	width, flags := ComputeProperties4(expression, openParen, arguments, closeParen)
	return &CallExpression{
		Base: Base{
			kind:      syntax.CallExpression,
			fullWidth: width,
			flags:     flags,
		},
		expression: expression,
		openParen:  openParen,
		arguments:  arguments,
		closeParen: closeParen,
	}
}

func (c *CallExpression) Expression() Node   { return c.expression }
func (c *CallExpression) OpenParen() *Token  { return c.openParen }
func (c *CallExpression) Arguments() Node    { return c.arguments }
func (c *CallExpression) CloseParen() *Token { return c.closeParen }

func (c *CallExpression) IsNil() bool { return c == nil }

func (c *CallExpression) SlotCount() int { return 4 }
func (c *CallExpression) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(c.expression)
	case 1:
		return SafeNode(c.openParen)
	case 2:
		return SafeNode(c.arguments)
	case 3:
		return SafeNode(c.closeParen)
	default:
		return nil
	}
}

type LiteralExpression struct {
	Base
	token *Token
}

func NewLiteralExpression(token *Token) *LiteralExpression {
	width, flags := ComputeProperties1(token)
	return &LiteralExpression{
		Base: Base{
			kind:      syntax.LiteralExpression,
			fullWidth: width,
			flags:     flags,
		},
		token: token,
	}
}

func (l *LiteralExpression) IsNil() bool { return l == nil }

func (l *LiteralExpression) Token() *Token  { return l.token }
func (l *LiteralExpression) SlotCount() int { return 1 }
func (l *LiteralExpression) Slot(index int) Node {
	if index == 0 {
		return SafeNode(l.token)
	}
	return nil
}

type MemberAccessExpression struct {
	Base
	expression Node
	dot        *Token
	name       *Token
}

func NewMemberAccessExpression(expression Node, dot, name *Token) *MemberAccessExpression {
	width, flags := ComputeProperties3(expression, dot, name)
	return &MemberAccessExpression{
		Base: Base{
			kind:      syntax.MemberAccessExpression,
			fullWidth: width,
			flags:     flags,
		},
		expression: expression,
		dot:        dot,
		name:       name,
	}
}

func (m *MemberAccessExpression) Expression() Node { return m.expression }
func (m *MemberAccessExpression) Dot() *Token      { return m.dot }
func (m *MemberAccessExpression) Name() *Token     { return m.name }

func (m *MemberAccessExpression) IsNil() bool { return m == nil }

func (m *MemberAccessExpression) SlotCount() int { return 3 }
func (m *MemberAccessExpression) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(m.expression)
	case 1:
		return SafeNode(m.dot)
	case 2:
		return SafeNode(m.name)
	default:
		return nil
	}
}

type NameExpression struct {
	Base
	identifier *Token
}

func NewNameExpression(identifier *Token) *NameExpression {
	width, flags := ComputeProperties1(identifier)
	return &NameExpression{
		Base: Base{
			kind:      syntax.NameExpression,
			fullWidth: width,
			flags:     flags,
		},
		identifier: identifier,
	}
}

func (n *NameExpression) IsNil() bool { return n == nil }

func (n *NameExpression) Identifier() *Token { return n.identifier }
func (n *NameExpression) SlotCount() int     { return 1 }
func (n *NameExpression) Slot(index int) Node {
	if index == 0 {
		return SafeNode(n.identifier)
	}
	return nil
}

type ParenthesizedExpression struct {
	Base
	openParen  *Token
	expression Node
	closeParen *Token
}

func NewParenthesizedExpression(openParen *Token, expression Node, closeParen *Token) *ParenthesizedExpression {
	width, flags := ComputeProperties3(openParen, expression, closeParen)
	return &ParenthesizedExpression{
		Base: Base{
			kind:      syntax.ParenthesizedExpression,
			fullWidth: width,
			flags:     flags,
		},
		openParen:  openParen,
		expression: expression,
		closeParen: closeParen,
	}
}

func (p *ParenthesizedExpression) OpenParen() *Token  { return p.openParen }
func (p *ParenthesizedExpression) Expression() Node   { return p.expression }
func (p *ParenthesizedExpression) CloseParen() *Token { return p.closeParen }

func (p *ParenthesizedExpression) IsNil() bool { return p == nil }

func (p *ParenthesizedExpression) SlotCount() int { return 3 }
func (p *ParenthesizedExpression) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(p.openParen)
	case 1:
		return SafeNode(p.expression)
	case 2:
		return SafeNode(p.closeParen)
	default:
		return nil
	}
}

type UnaryExpression struct {
	Base
	operator *Token
	operand  Node
}

func NewUnaryExpression(operator *Token, operand Node) *UnaryExpression {
	width, flags := ComputeProperties2(operator, operand)
	return &UnaryExpression{
		Base: Base{
			kind:      syntax.UnaryExpression,
			fullWidth: width,
			flags:     flags,
		},
		operator: operator,
		operand:  operand,
	}
}

func (u *UnaryExpression) Operator() *Token { return u.operator }
func (u *UnaryExpression) Operand() Node    { return u.operand }

func (u *UnaryExpression) IsNil() bool { return u == nil }

func (u *UnaryExpression) SlotCount() int { return 2 }
func (u *UnaryExpression) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(u.operator)
	case 1:
		return SafeNode(u.operand)
	default:
		return nil
	}
}
