package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func optimizeExpr(expr ast.Expr) ast.Expr {
	if expr == nil {
		return nil
	}

	switch n := expr.(type) {
	case *ast.BinaryExpr:
		// Bottom-up: optimize children first
		n.Left = optimizeExpr(n.Left)
		n.Right = optimizeExpr(n.Right)

		// Try constant folding
		if folded := foldBinaryExpr(n.Left, n.Operator, n.Right); folded != nil {
			return optimizeExpr(folded)
		}

		// Try algebraic/boolean simplifications
		if simplified := simplifyBinary(n); simplified != nil {
			return optimizeExpr(simplified)
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
		exprs[i] = optimizeExpr(exprs[i])
	}
}
