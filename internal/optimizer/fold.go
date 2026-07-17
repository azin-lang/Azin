package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func foldBinaryExpr(left ast.Expr, op token.Token, right ast.Expr) ast.Expr {
	switch l := left.(type) {
	case *ast.IntegerLiteral:
		switch r := right.(type) {
		case *ast.IntegerLiteral:
			return foldIntegerInteger(l, op, r)

		case *ast.FloatLiteral:
			return foldFloatFloat(
				floatLit(float64(l.Value)),
				op,
				r,
			)

		case *ast.CharacterLiteral:
			return foldIntegerInteger(
				l,
				op,
				intLit(int64(r.Value)),
			)
		}

	case *ast.FloatLiteral:
		switch r := right.(type) {
		case *ast.IntegerLiteral:
			return foldFloatFloat(
				l,
				op,
				floatLit(float64(r.Value)),
			)

		case *ast.FloatLiteral:
			return foldFloatFloat(l, op, r)
		}

	case *ast.BooleanLiteral:
		if r, ok := right.(*ast.BooleanLiteral); ok {
			return foldBooleanBoolean(l, op, r)
		}

	case *ast.CharacterLiteral:
		switch r := right.(type) {
		case *ast.CharacterLiteral:
			return foldIntegerInteger(
				charAsInt(l),
				op,
				charAsInt(r),
			)

		case *ast.IntegerLiteral:
			return foldIntegerInteger(
				charAsInt(l),
				op,
				r,
			)
		}
	}

	return nil
}
