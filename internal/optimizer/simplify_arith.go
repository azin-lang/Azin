package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func simplifyArithmetic(n *ast.BinaryExpr) ast.Expr {
	leftPure := isPure(n.Left)

	switch n.Operator.Kind {

	case token.Plus, token.Minus:
		if isZero(n.Right) {
			return n.Left
		}
		if n.Operator.Kind == token.Minus && leftPure && exprEqual(n.Left, n.Right) {
			return intLit(0)
		}

	case token.Star:
		// x * 1 == x
		if isOne(n.Right) {
			return n.Left
		}

		// x * 0 == 0 iff x is pure
		if isZero(n.Right) && leftPure {
			return n.Right
		}

		// x * 2 == x + x (if pure)
		if leftPure {
			if r, ok := n.Right.(*ast.IntegerLiteral); ok && r.Value == 2 {
				return &ast.BinaryExpr{
					Left:     n.Left,
					Operator: token.Token{Kind: token.Plus},
					Right:    n.Left,
				}
			}
		}

		// TODO: add x * 2^k == x << k (only if pure)
		// Maybe also check if k > width of given literal (if known) and return 0?

	case token.Slash:
		// x / 1
		if isOne(n.Right) {
			return n.Left
		}

		// TODO: add x / 2^k => x >> k

		// x / x == 1 only if x != 0 and x is pure
		if leftPure && exprEqual(n.Left, n.Right) && !isZero(n.Left) {
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
