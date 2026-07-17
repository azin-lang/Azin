package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func isTrue(expr ast.Expr) bool {
	b, ok := expr.(*ast.BooleanLiteral)
	return ok && b.Value
}

func isFalse(expr ast.Expr) bool {
	b, ok := expr.(*ast.BooleanLiteral)
	return ok && !b.Value
}

// TODO: move this out of optimizer and implement an `Equal(Expr) bool` inside the Expr interface

func exprEqual(left, right ast.Expr) bool {
	switch l := left.(type) {
	case *ast.IntegerLiteral:
		r, ok := right.(*ast.IntegerLiteral)
		return ok && l.Value == r.Value

	case *ast.FloatLiteral:
		r, ok := right.(*ast.FloatLiteral)
		return ok && l.Value == r.Value

	case *ast.BooleanLiteral:
		r, ok := right.(*ast.BooleanLiteral)
		return ok && l.Value == r.Value

	case *ast.CharacterLiteral:
		r, ok := right.(*ast.CharacterLiteral)
		return ok && l.Value == r.Value

	case *ast.Identifier:
		r, ok := right.(*ast.Identifier)
		return ok && l.Value == r.Value

	case *ast.MemberExpr:
		r, ok := right.(*ast.MemberExpr)
		return ok &&
			exprEqual(l.Object, r.Object) &&
			exprEqual(l.Property, r.Property)

	case *ast.BinaryExpr:
		r, ok := right.(*ast.BinaryExpr)
		return ok &&
			l.Operator.Kind == r.Operator.Kind &&
			exprEqual(l.Left, r.Left) &&
			exprEqual(l.Right, r.Right)

	case *ast.CallExpr:
		r, ok := right.(*ast.CallExpr)
		if !ok {
			return false
		}

		if !exprEqual(l.Callee, r.Callee) {
			return false
		}

		if len(l.Args) != len(r.Args) {
			return false
		}

		for i := range l.Args {
			if !exprEqual(l.Args[i], r.Args[i]) {
				return false
			}
		}

		return true

	default:
		return false
	}
}
