package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type LiteralExpression struct {
	Base

	token *Token
}

func NewLiteralExpression(token *Token) *LiteralExpression {
	var width uint32
	if token != nil {
		width += token.FullWidth()
	}

	return &LiteralExpression{
		Base: Base{
			kind:      syntax.LiteralExpression,
			fullWidth: width,
		},
		token: token,
	}
}

func (l *LiteralExpression) Token() *Token { return l.token }

func (l *LiteralExpression) SlotCount() int { return 1 }

func (l *LiteralExpression) Slot(index int) Node {
	if index == 0 {
		return l.token
	}
	return nil
}
