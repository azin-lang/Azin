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

		// x * 2^k == x << k (if pure)
		if leftPure {
			if k, ok := isPowerOfTwo(n.Right); ok {
				return &ast.BinaryExpr{
					Left:     n.Left,
					Operator: token.Token{Kind: token.LessLess},
					Right:    intLit(k),
				}
			}
		}

	case token.Slash:
		// x / 1
		if isOne(n.Right) {
			return n.Left
		}

		// x / 2^k == x >> k (if pure)
		if leftPure {
			if k, ok := isPowerOfTwo(n.Right); ok {
				return &ast.BinaryExpr{
					Left:     n.Left,
					Operator: token.Token{Kind: token.GreaterGreater},
					Right:    intLit(k),
				}
			}
		}

		// x / x == 1 only if x != 0 and x is pure
		if leftPure && exprEqual(n.Left, n.Right) && !isZero(n.Left) {
			return intLit(1)
		}

	case token.Modulo:
		// x % 1
		if isOne(n.Right) {
			return intLit(0)
		}

		// x % 2^k == x & (2^k - 1) (if pure)
		if leftPure {
			if k, ok := isPowerOfTwo(n.Right); ok {
				return &ast.BinaryExpr{
					Left:     n.Left,
					Operator: token.Token{Kind: token.Ampersand},
					Right:    intLit((1 << k) - 1),
				}
			}
		}
	}

	return nil
}
