package opt

import (
	"github.com/azin-lang/Azin/pkg/ast"
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func foldInteger(left *ast.IntegerLiteral, op token2.Token, right *ast.IntegerLiteral) ast.Expr {
	switch op.Kind {
	case token2.Plus:
		return intLit(left.Value + right.Value)
	case token2.Minus:
		return intLit(left.Value - right.Value)
	case token2.Star:
		return intLit(left.Value * right.Value)
	case token2.Slash:
		if right.Value == 0 {
			return nil
		}
		return intLit(left.Value / right.Value)
	case token2.Modulo:
		if right.Value == 0 {
			return nil
		}
		return intLit(left.Value % right.Value)
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
