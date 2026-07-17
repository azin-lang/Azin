package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func reassociateBinary(n *ast.BinaryExpr) ast.Expr {
	left, ok := n.Left.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	if left.Operator.Kind != n.Operator.Kind {
		return nil
	}

	switch n.Operator.Kind {
	case token.Plus:
		l, ok1 := left.Right.(*ast.IntegerLiteral)
		r, ok2 := n.Right.(*ast.IntegerLiteral)

		if ok1 && ok2 {
			left.Right = intLit(l.Value + r.Value)
			return left
		}
	case token.Star:
		l, ok1 := left.Right.(*ast.IntegerLiteral)
		r, ok2 := n.Right.(*ast.IntegerLiteral)

		if ok1 && ok2 {
			left.Right = intLit(l.Value * r.Value)
			return left
		}
	}
	return nil
}

func canonicalizeBinary(n *ast.BinaryExpr) ast.Expr {
	switch n.Operator.Kind {

	case token.Plus,
		token.Star,
		token.EqualEqual,
		token.BangEqual:

		if isConstant(n.Left) && !isConstant(n.Right) {
			n.Left, n.Right = n.Right, n.Left
			return n
		}
	}

	return nil
}
