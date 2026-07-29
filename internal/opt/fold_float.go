package opt

import (
	"github.com/azin-lang/Azin/pkg/ast"
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func foldFloat(left *ast.FloatLiteral, op token2.Token, right *ast.FloatLiteral) ast.Expr {
	switch op.Kind {
	case token2.Plus:
		return floatLit(left.Value + right.Value)
	case token2.Minus:
		return floatLit(left.Value - right.Value)
	case token2.Star:
		return floatLit(left.Value * right.Value)
	case token2.Slash:
		if right.Value == 0 {
			return nil
		}
		return floatLit(left.Value / right.Value)
	case token2.EqualEqual:
		return boolLit(left.Value == right.Value)
	case token2.BangEqual:
		return boolLit(left.Value != right.Value)
	case token2.Less:
		return boolLit(left.Value < right.Value)
	case token2.LessEqual:
		return boolLit(left.Value <= right.Value)
	case token2.Greater:
		return boolLit(left.Value > right.Value)
	case token2.GreaterEqual:
		return boolLit(left.Value >= right.Value)
	default:
		return nil
	}
}
