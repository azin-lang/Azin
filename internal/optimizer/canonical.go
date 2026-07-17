package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

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
