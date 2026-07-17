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

	case token.Minus:
		// x - 0
		if isZero(n.Right) {
			return n.Left
		}

		// x - x
		if isPure(n.Left) && exprEqual(n.Left, n.Right) {
			return intLit(0)
		}

	case token.Star:
		// x * 2 == x + x (if pure)
		if isPure(n.Left) {
			r, ok := n.Right.(*ast.IntegerLiteral)
			if !ok || r.Value != 2 || !isPure(n.Left) {
				return nil
			}

			return &ast.BinaryExpr{
				Left:     n.Left,
				Operator: token.Token{Kind: token.Plus},
				Right:    n.Left,
			}
		}

		// TODO: add x * 2^k == x << k (only if pure)
		// Maybe also check if k > width of given literal (if known) and return 0?

		// x * 1
		if isOne(n.Right) {
			return n.Left
		}

		// x * 0 => 0 (if pure)
		if isZero(n.Right) && isPure(n.Left) {
			return n.Right
		}

	case token.Slash:
		// x / 1
		if isOne(n.Right) {
			return n.Left
		}

		// TODO: add x / 2^k => x >> k

		// x / x == 1 only if x != 0 and x is pure
		if isPure(n.Left) && exprEqual(n.Left, n.Right) && !isZero(n.Left) {
			return intLit(1)
		}

	case token.Modulo:
		// x % 1
		if isOne(n.Right) {
			return intLit(0)
		}

		// TODO: add x % 2^k => x & (2^k - 1)
	}

	return nil
}
