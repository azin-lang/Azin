package c

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/token"
)

func emitOperator(kind token.Kind) string {
	switch kind {
	case token.Plus:
		return "+"
	case token.Minus:
		return "-"
	case token.Star:
		return "*"
	case token.Slash:
		return "/"
	case token.EqualEqual:
		return "=="
	case token.BangEqual:
		return "!="
	case token.Less:
		return "<"
	case token.LessEqual:
		return "<="
	case token.Greater:
		return ">"
	case token.GreaterEqual:
		return ">="
	default:
		panic(fmt.Sprintf("unsupported operator %v", kind))
	}
}
