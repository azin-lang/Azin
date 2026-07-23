package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func (o *Optimizer) optimizeExpr(expr ast.Expr) ast.Expr {
	if expr == nil {
		return nil
	}

	switch n := expr.(type) {
	case *ast.Identifier:
		if val, ok := o.currentScope.GetConstant(n.Value); ok {
			return cloneConstant(val)
		}

	case *ast.BinaryExpr:
		// Bottom-up: optimize children first
		n.Left = o.optimizeExpr(n.Left)
		n.Right = o.optimizeExpr(n.Right)

		// Try constant folding
		if folded := foldBinaryExpr(n.Left, n.Operator, n.Right); folded != nil {
			return o.optimizeExpr(folded)
		}

		// Try algebraic/boolean simplifications
		if simplified := simplifyBinary(n); simplified != nil {
			return o.optimizeExpr(simplified)
		}

	case *ast.MemberExpr:
		n.Object = o.optimizeExpr(n.Object)

	case *ast.CallExpr:
		n.Callee = o.optimizeExpr(n.Callee)
		o.optimizeExprs(n.Args)
	}

	return expr
}

func (o *Optimizer) optimizeExprs(exprs []ast.Expr) {
	for i := range exprs {
		exprs[i] = o.optimizeExpr(exprs[i])
	}
}
