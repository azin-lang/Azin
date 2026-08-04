package green

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

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

func (l *LiteralExpression) Token() *Token  { return l.token }
func (l *LiteralExpression) SlotCount() int { return 1 }
func (l *LiteralExpression) Slot(index int) Node {
	if index == 0 && l.token != nil {
		return l.token
	}
	return nil
}
