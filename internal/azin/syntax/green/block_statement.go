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

func (b *BlockStatement) SlotCount() int { return 3 }
func (b *BlockStatement) Slot(index int) Node {
	switch index {
	case 0:
		if b.startToken == nil {
			return nil
		}
		return b.startToken
	case 1:
		if b.statements == nil {
			return nil
		}
		return b.statements
	case 2:
		if b.endToken == nil {
			return nil
		}
		return b.endToken
	default:
		return nil
	}
}
