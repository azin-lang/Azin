package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func foldInteger(left *ast.IntegerLiteral, op token.Token, right *ast.IntegerLiteral) ast.Expr {
	switch op.Kind {
	case token.Plus:
		return intLit(left.Value + right.Value)
	case token.Minus:
		return intLit(left.Value - right.Value)
	case token.Star:
		return intLit(left.Value * right.Value)
	case token.Slash:
		if right.Value == 0 {
			return nil
		}
		return intLit(left.Value / right.Value)
	case token.Modulo:
		if right.Value == 0 {
			return nil
		}
		return intLit(left.Value % right.Value)
	case token.EqualEqual:
		return boolLit(left.Value == right.Value)
	case token.BangEqual:
		return boolLit(left.Value != right.Value)
	case token.Less:
		return boolLit(left.Value < right.Value)
	case token.LessEqual:
		return boolLit(left.Value <= right.Value)
	case token.Greater:
		return boolLit(left.Value > right.Value)
	case token.GreaterEqual:
		return boolLit(left.Value >= right.Value)
	default:
		return nil
	}
}
