package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type BlockStatement struct {
	Base
	startToken *Token
	statements Node
	endToken   *Token
}

func NewBlockStatement(startToken *Token, statements Node, endToken *Token) *BlockStatement {
	width, flags := ComputeProperties3(startToken, statements, endToken)
	return &BlockStatement{
		Base: Base{
			kind:      syntax.BlockStatement,
			fullWidth: width,
			flags:     flags,
		},
		startToken: startToken,
		statements: statements,
		endToken:   endToken,
	}
}

func (b *BlockStatement) StartToken() *Token { return b.startToken }
func (b *BlockStatement) Statements() Node   { return b.statements }
func (b *BlockStatement) EndToken() *Token   { return b.endToken }

func (b *BlockStatement) IsNil() bool { return b == nil }

func (b *BlockStatement) SlotCount() int { return 3 }
func (b *BlockStatement) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(b.startToken)
	case 1:
		return SafeNode(b.statements)
	case 2:
		return SafeNode(b.endToken)
	default:
		return nil
	}
}

type IfStatement struct {
	Base
	ifKeyword   *Token
	condition   Node
	thenBranch  Node
	elseKeyword *Token
	elseBranch  Node
}

func NewIfStatement(ifKeyword *Token, condition Node, thenBranch Node, elseKeyword *Token, elseBranch Node) *IfStatement {
	width, flags := ComputeProperties5(ifKeyword, condition, thenBranch, elseKeyword, elseBranch)
	return &IfStatement{
		Base: Base{
			kind:      syntax.IfStatement,
			fullWidth: width,
			flags:     flags,
		},
		ifKeyword:   ifKeyword,
		condition:   condition,
		thenBranch:  thenBranch,
		elseKeyword: elseKeyword,
		elseBranch:  elseBranch,
	}
}

func (i *IfStatement) IfKeyword() *Token   { return i.ifKeyword }
func (i *IfStatement) Condition() Node     { return i.condition }
func (i *IfStatement) ThenBranch() Node    { return i.thenBranch }
func (i *IfStatement) ElseKeyword() *Token { return i.elseKeyword }
func (i *IfStatement) ElseBranch() Node    { return i.elseBranch }

func (i *IfStatement) IsNil() bool { return i == nil }

func (i *IfStatement) SlotCount() int { return 5 }
func (i *IfStatement) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(i.ifKeyword)
	case 1:
		return SafeNode(i.condition)
	case 2:
		return SafeNode(i.thenBranch)
	case 3:
		return SafeNode(i.elseKeyword)
	case 4:
		return SafeNode(i.elseBranch)
	default:
		return nil
	}
}

type ReturnStatement struct {
	Base
	returnKeyword *Token
	expression    Node
}

func NewReturnStatement(returnKeyword *Token, expression Node) *ReturnStatement {
	width, flags := ComputeProperties2(returnKeyword, expression)
	return &ReturnStatement{
		Base: Base{
			kind:      syntax.ReturnStatement,
			fullWidth: width,
			flags:     flags,
		},
		returnKeyword: returnKeyword,
		expression:    expression,
	}
}

func (r *ReturnStatement) ReturnKeyword() *Token { return r.returnKeyword }
func (r *ReturnStatement) Expression() Node      { return r.expression }

func (r *ReturnStatement) IsNil() bool { return r == nil }

func (r *ReturnStatement) SlotCount() int { return 2 }
func (r *ReturnStatement) Slot(index int) Node {
	switch index {
	case 0:
		return SafeNode(r.returnKeyword)
	case 1:
		return SafeNode(r.expression)
	default:
		return nil
	}
}
