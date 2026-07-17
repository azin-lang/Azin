package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func simplifyBoolean(n *ast.BinaryExpr) ast.Expr {
	switch n.Operator.Kind {

	case token.LogicalAnd:
		// true && x -> x
		if isTrue(n.Left) {
			return n.Right
		}

		// x && true -> x
		if isTrue(n.Right) {
			return n.Left
		}

		// false && x -> false
		if isFalse(n.Left) {
			return n.Left
		}

		// x && false -> false
		if isFalse(n.Right) {
			return n.Right
		}

		// x && x -> x
		if exprEqual(n.Left, n.Right) {
			return n.Left
		}

	case token.LogicalOr:
		// false || x -> x
		if isFalse(n.Left) {
			return n.Right
		}

		// x || false -> x
		if isFalse(n.Right) {
			return n.Left
		}

		// true || x -> true
		if isTrue(n.Left) {
			return n.Left
		}

		// x || true -> true
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

		// x == true
		if isTrue(n.Right) {
			return n.Left
		}

		// true == x
		if isTrue(n.Left) {
			return n.Right
		}
	}

	return nil
}
