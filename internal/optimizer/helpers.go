package optimizer

import (
	"math/bits"

	"github.com/azin-lang/Azin/internal/ast"
)

func isTrue(expr ast.Expr) bool {
	b, ok := expr.(*ast.BooleanLiteral)
	return ok && b.Value
}

func isFalse(expr ast.Expr) bool {
	b, ok := expr.(*ast.BooleanLiteral)
	return ok && !b.Value
}

func exprEqual(left, right ast.Expr) bool {
	// Fast path: identical pointers
	if left == right {
		return true
	}

	if left == nil || right == nil {
		return false
	}

	return left.Equals(right)
}

func isZero(expr ast.Expr) bool {
	switch n := expr.(type) {
	case *ast.IntegerLiteral:
		return n.Value == 0
	case *ast.FloatLiteral:
		return n.Value == 0
	case *ast.CharacterLiteral:
		return n.Value == 0
	}
	return false
}

func isPowerOfTwo(expr ast.Expr) (int64, bool) {
	n, ok := expr.(*ast.IntegerLiteral)
	if !ok || n.Value <= 0 {
		return 0, false
	}

	v := uint64(n.Value)
	if v&(v-1) != 0 {
		return 0, false
	}

	return int64(bits.TrailingZeros64(v)), true
}

// isNonNegative returns true if the expression is definitely >= 0.
// The optimizer has no type information, so this only applies to literals.
func isNonNegative(expr ast.Expr) bool {
	switch n := expr.(type) {
	case *ast.IntegerLiteral:
		return n.Value >= 0
	case *ast.CharacterLiteral:
		return n.Value >= 0
	case *ast.BooleanLiteral:
		return true
	}
	return false
}

func isOne(expr ast.Expr) bool {
	switch n := expr.(type) {
	case *ast.IntegerLiteral:
		return n.Value == 1
	case *ast.FloatLiteral:
		return n.Value == 1
	case *ast.CharacterLiteral:
		return n.Value == 1
	}
	return false
}

func isConstant(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.IntegerLiteral, *ast.FloatLiteral, *ast.BooleanLiteral, *ast.CharacterLiteral:
		return true
	}
	return false
}

// isNotFloat returns true if the expression is definitely not a float type.
// The optimizer has no type information, so we can only determine this for literals.
func isNotFloat(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.IntegerLiteral, *ast.BooleanLiteral, *ast.CharacterLiteral, *ast.StringLiteral:
		return true
	case *ast.FloatLiteral:
		return false
	default:
		// Identifiers, binary expressions, call expressions, etc. could be float.
		return false
	}
}

// isPure returns true if evaluating the expression has no side effects.
func isPure(expr ast.Expr) bool {
	switch n := expr.(type) {
	case *ast.IntegerLiteral, *ast.FloatLiteral, *ast.BooleanLiteral, *ast.CharacterLiteral, *ast.Identifier:
		return true
	case *ast.MemberExpr:
		return isPure(n.Object)
	case *ast.BinaryExpr:
		return isPure(n.Left) && isPure(n.Right)
	default:
		return false
	}
}

func isCopyable(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.IntegerLiteral, *ast.FloatLiteral, *ast.BooleanLiteral, *ast.CharacterLiteral, *ast.StringLiteral:
		return true
	case *ast.Identifier:
		// An identifier can be propagated as long as it's pure/safe.
		return true
	}
	return false
}

func cloneValue(expr ast.Expr) ast.Expr {
	switch n := expr.(type) {
	case *ast.IntegerLiteral:
		return &ast.IntegerLiteral{Token: n.Token, Value: n.Value}
	case *ast.FloatLiteral:
		return &ast.FloatLiteral{Token: n.Token, Value: n.Value}
	case *ast.BooleanLiteral:
		return &ast.BooleanLiteral{Token: n.Token, Value: n.Value}
	case *ast.CharacterLiteral:
		return &ast.CharacterLiteral{Token: n.Token, Value: n.Value}
	case *ast.StringLiteral:
		return &ast.StringLiteral{Token: n.Token, Value: n.Value}
	case *ast.Identifier:
		return &ast.Identifier{Token: n.Token, Value: n.Value}
	}
	return expr
}
