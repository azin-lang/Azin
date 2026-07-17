package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func simplifyArithmetic(n *ast.BinaryExpr) ast.Expr {
	switch n.Operator.Kind {

	case token.Plus:
		// x + 0
		if isZero(n.Right) {
			return n.Left
		}

		// 0 + x
		if isZero(n.Left) {
			return n.Right
		}

	case token.Minus:
		// x - 0
		if isZero(n.Right) {
			return n.Left
		}

		// x - x
		if exprEqual(n.Left, n.Right) {
			return intLit(0)
		}

	case token.Star:
		// x * 1
		if isOne(n.Right) {
			return n.Left
		}

		// 1 * x
		if isOne(n.Left) {
			return n.Right
		}

		// x * 0
		if isZero(n.Right) {
			return n.Right
		}

		// 0 * x
		if isZero(n.Left) {
			return n.Left
		}

	case token.Slash:
		// x / 1
		if isOne(n.Right) {
			return n.Left
		}

	case token.Modulo:
		// x % 1
		if isOne(n.Right) {
			return intLit(0)
		}
	}

	return nil
}
