package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func optimizeExpr(expr ast.Expr) ast.Expr {
	for {
		next := optimizeExprOnce(expr)
		if exprEqual(expr, next) {
			return next
		}
		expr = next
	}
}

func optimizeExprOnce(expr ast.Expr) ast.Expr {
	switch n := expr.(type) {

	case *ast.BinaryExpr:
		n.Left = optimizeExpr(n.Left)
		n.Right = optimizeExpr(n.Right)

		if folded := foldBinaryExpr(n.Left, n.Operator, n.Right); folded != nil {
			return folded
		}

		if simplified := simplifyBinary(n); simplified != nil {
			return simplified
		}

	case *ast.MemberExpr:
		n.Object = optimizeExpr(n.Object)

	case *ast.CallExpr:
		n.Callee = optimizeExpr(n.Callee)
		optimizeExprs(n.Args)
	}

	return expr
}

func optimizeExprs(exprs []ast.Expr) {
	for i := range exprs {
		exprs[i] = optimizeExprOnce(exprs[i])
	}
}
