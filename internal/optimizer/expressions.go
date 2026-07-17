package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func optimizeExpr(expr ast.Expr) ast.Expr {
	switch n := expr.(type) {
	case *ast.BinaryExpr:
		n.Left = optimizeExpr(n.Left)
		n.Right = optimizeExpr(n.Right)

		if expr := simplifyBinary(n); expr != nil {
			return expr
		}

		if folded := foldBinaryExpr(n.Left, n.Operator, n.Right); folded != nil {
			return folded
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

func simplifyBinary(n *ast.BinaryExpr) ast.Expr {
	switch n.Operator.Kind {

	case token.LogicalAnd:
		if isTrue(n.Left) {
			return n.Right
		}
		if isTrue(n.Right) {
			return n.Left
		}

		if isFalse(n.Left) {
			return n.Left
		}
		if isFalse(n.Right) {
			return n.Right
		}

		// x && x -> x
		if exprEqual(n.Left, n.Right) {
			return n.Left
		}

	case token.LogicalOr:
		if isFalse(n.Left) {
			return n.Right
		}
		if isFalse(n.Right) {
			return n.Left
		}

		if isTrue(n.Left) {
			return n.Left
		}
		if isTrue(n.Right) {
			return n.Right
		}

		// x || x -> x
		if exprEqual(n.Left, n.Right) {
			return n.Left
		}

	case token.EqualEqual:
		// x == x
		if exprEqual(n.Left, n.Right) {
			return boolLit(true)
		}

	case token.BangEqual:
		// x != x
		if exprEqual(n.Left, n.Right) {
			return boolLit(false)
		}
	}

	return nil
}
