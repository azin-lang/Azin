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
			return foldInteger(l, op, r)

		case *ast.FloatLiteral:
			return foldFloat(
				floatLit(float64(l.Value)),
				op,
				r,
			)

		case *ast.CharacterLiteral:
			return foldInteger(
				l,
				op,
				intLit(int64(r.Value)),
			)
		}

	case *ast.FloatLiteral:
		switch r := right.(type) {
		case *ast.IntegerLiteral:
			return foldFloat(
				l,
				op,
				floatLit(float64(r.Value)),
			)

		case *ast.FloatLiteral:
			return foldFloat(l, op, r)
		}

	case *ast.BooleanLiteral:
		if r, ok := right.(*ast.BooleanLiteral); ok {
			return foldBoolean(l, op, r)
		}

	case *ast.CharacterLiteral:
		switch r := right.(type) {
		case *ast.CharacterLiteral:
			return foldInteger(
				charAsInt(l),
				op,
				charAsInt(r),
			)

		case *ast.IntegerLiteral:
			return foldInteger(
				charAsInt(l),
				op,
				r,
			)
		}
	}

	return nil
}
