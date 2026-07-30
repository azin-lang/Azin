package optimize

import (
	"github.com/azin-lang/Azin/pkg/ast"
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func simplifyArithmetic(n *ast.BinaryExpr) ast.Expr {
	leftPure := isPure(n.Left)
	//nolint:exhaustive
	switch n.Operator.Kind {

	case token2.Plus, token2.Minus:
		if isZero(n.Right) {
			return n.Left
		}
		if n.Operator.Kind == token2.Minus && leftPure && exprEqual(n.Left, n.Right) {
			return intLit(0)
		}

	case token2.Star:
		// x * 1 == x
		if isOne(n.Right) {
			return n.Left
		}

		// x * 0 == 0 iff x is pure
		if isZero(n.Right) && leftPure {
			return n.Right
		}

		// x * 2^k == x << k (if pure)
		if leftPure {
			if k, ok := isPowerOfTwo(n.Right); ok {
				return &ast.BinaryExpr{
					Left:     n.Left,
					Operator: token2.Token{Kind: token2.LessLess},
					Right:    intLit(k),
				}
			}
		}

	case token2.Slash:
		// x / 1
		if isOne(n.Right) {
			return n.Left
		}

		// x / 2^k == x >> k (if pure and non-negative)
		if leftPure && isNonNegative(n.Left) {
			if k, ok := isPowerOfTwo(n.Right); ok {
				return &ast.BinaryExpr{
					Left:     n.Left,
					Operator: token2.Token{Kind: token2.GreaterGreater},
					Right:    intLit(k),
				}
			}
		}

		// x / x == 1 only if x != 0 and x is pure
		if leftPure && exprEqual(n.Left, n.Right) && !isZero(n.Left) {
			return intLit(1)
		}

	case token2.Modulo:
		// x % 1
		if isOne(n.Right) {
			return intLit(0)
		}

		// x % 2^k == x & (2^k - 1) (if pure and non-negative)
		if leftPure && isNonNegative(n.Left) {
			if k, ok := isPowerOfTwo(n.Right); ok {
				return &ast.BinaryExpr{
					Left:     n.Left,
					Operator: token2.Token{Kind: token2.Ampersand},
					Right:    intLit((1 << k) - 1),
				}
			}
		}
	}

	return nil
}
