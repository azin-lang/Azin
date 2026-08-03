package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type IfStatement struct {
	Base
	ifKeyword   *Token
	condition   Node
	thenBranch  Node
	elseKeyword *Token
	elseBranch  Node
}

func NewIfStatement(ifKeyword *Token, condition Node, thenBranch Node, elseKeyword *Token, elseBranch Node) *IfStatement {
	var width uint32
	if ifKeyword != nil {
		width += ifKeyword.FullWidth()
	}
	if condition != nil {
		width += condition.FullWidth()
	}
	if thenBranch != nil {
		width += thenBranch.FullWidth()
	}
	if elseKeyword != nil {
		width += elseKeyword.FullWidth()
	}
	if elseBranch != nil {
		width += elseBranch.FullWidth()
	}

	return &IfStatement{
		Base: Base{
			kind:      syntax.IfStatement,
			fullWidth: width,
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

func (i *IfStatement) SlotCount() int { return 5 }

func (i *IfStatement) Slot(index int) Node {
	switch index {
	case 0:
		return i.ifKeyword
	case 1:
		return i.condition
	case 2:
		return i.thenBranch
	case 3:
		return i.elseKeyword
	case 4:
		return i.elseBranch
	default:
		return nil
	}
}
