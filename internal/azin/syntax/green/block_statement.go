package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type BlockStatement struct {
	Base
	openBrace  *Token
	statements Node // SyntaxList of statements
	closeBrace *Token
}

func NewBlockStatement(openBrace *Token, statements Node, closeBrace *Token) *BlockStatement {
	var width uint32
	if openBrace != nil {
		width += openBrace.FullWidth()
	}
	if statements != nil {
		width += statements.FullWidth()
	}
	if closeBrace != nil {
		width += closeBrace.FullWidth()
	}

	return &BlockStatement{
		Base: Base{
			kind:      syntax.BlockStatement,
			fullWidth: width,
		},
		openBrace:  openBrace,
		statements: statements,
		closeBrace: closeBrace,
	}
}

func (b *BlockStatement) OpenBrace() *Token  { return b.openBrace }
func (b *BlockStatement) Statements() Node   { return b.statements }
func (b *BlockStatement) CloseBrace() *Token { return b.closeBrace }

func (b *BlockStatement) SlotCount() int { return 3 }

func (b *BlockStatement) Slot(index int) Node {
	switch index {
	case 0:
		return b.openBrace
	case 1:
		return b.statements
	case 2:
		return b.closeBrace
	default:
		return nil
	}
}
