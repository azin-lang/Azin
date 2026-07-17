package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func simplifyBoolean(n *ast.BinaryExpr) ast.Expr {
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
		if isFalse(n.Right) && isPure(n.Left) {
			return n.Right
		}

		// x && x == x
		if isPure(n.Left) && exprEqual(n.Left, n.Right) {
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
		if isTrue(n.Right) && isPure(n.Left) {
			return n.Right
		}

		// x || x == x
		if isPure(n.Left) && exprEqual(n.Left, n.Right) {
			return n.Left
		}

	case token.EqualEqual:
		// x == x is true
		// TODO: handle NaN, since the identity for that
		if isPure(n.Left) && exprEqual(n.Left, n.Right) {
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
