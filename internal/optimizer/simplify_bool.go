package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func simplifyBoolean(n *ast.BinaryExpr) ast.Expr {
	leftPure := isPure(n.Left)

	switch n.Operator.Kind {

	case token.LogicalAnd:
		// true && x == x
		if isTrue(n.Left) {
			return n.Right
		}

		// x && true == x
		if isTrue(n.Right) {
			return n.Left
		}

		// false && foo() safely short-circuits to false
		if isFalse(n.Left) {
			return n.Left
		}

		// foo() && false is only false if foo() is pure
		if isFalse(n.Right) && leftPure {
			return n.Right
		}

		// x && x == x
		if leftPure && exprEqual(n.Left, n.Right) {
			return n.Left
		}

	case token.LogicalOr:
		// false || x == x
		if isFalse(n.Left) {
			return n.Right
		}

		// x || false == x
		if isFalse(n.Right) {
			return n.Left
		}

		// true || foo() safely short-circuits to true
		if isTrue(n.Left) {
			return n.Left
		}

		// foo() || true is only true if foo() is pure
		if isTrue(n.Right) && leftPure {
			return n.Right
		}

		// x || x == x
		if leftPure && exprEqual(n.Left, n.Right) {
			return n.Left
		}

	case token.EqualEqual:
		// x == x is true (but not for floats, NaN != NaN)
		if leftPure && exprEqual(n.Left, n.Right) && isNotFloat(n.Left) {
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

	case token.Less, token.Greater:
		// x < x == false, x > x == false
		if leftPure && exprEqual(n.Left, n.Right) {
			return boolLit(false)
		}

	case token.LessEqual, token.GreaterEqual:
		// x <= x == true, x >= x == true (but not for floats, NaN <= NaN is false)
		if leftPure && exprEqual(n.Left, n.Right) && isNotFloat(n.Left) {
			return boolLit(true)
		}
	}

	return nil
}
