package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func simplifyBinary(n *ast.BinaryExpr) ast.Expr {
	if expr := simplifyBoolean(n); expr != nil {
		return expr
	}

	if expr := simplifyArithmetic(n); expr != nil {
		return expr
	}

	return nil
}
