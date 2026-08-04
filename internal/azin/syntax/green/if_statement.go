package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

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

func (i *IfStatement) SlotCount() int { return 5 }
func (i *IfStatement) Slot(index int) Node {
	switch index {
	case 0:
		if i.ifKeyword == nil {
			return nil
		}
		return i.ifKeyword
	case 1:
		if i.condition == nil {
			return nil
		}
		return i.condition
	case 2:
		if i.thenBranch == nil {
			return nil
		}
		return i.thenBranch
	case 3:
		if i.elseKeyword == nil {
			return nil
		}
		return i.elseKeyword
	case 4:
		if i.elseBranch == nil {
			return nil
		}
		return i.elseBranch
	default:
		return nil
	}
}
