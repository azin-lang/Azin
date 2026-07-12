package codegen

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

func emitType(name string) string {
	switch name {
	case "unit":
		return "void"

	case "int":
		return "int"

	case "float":
		return "float"

	case "char":
		return "char"

	case "string":
		return "const char *"

	default:
		panic("unknown type: " + name)
	}
}
