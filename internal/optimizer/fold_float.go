package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func foldFloatFloat(left *ast.FloatLiteral, op token.Token, right *ast.FloatLiteral) ast.Expr {
	switch op.Kind {
	case token.Plus:
		return floatLit(left.Value + right.Value)
	case token.Minus:
		return floatLit(left.Value - right.Value)
	case token.Star:
		return floatLit(left.Value * right.Value)
	case token.Slash:
		if right.Value == 0 {
			return nil
		}
		return floatLit(left.Value / right.Value)
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
