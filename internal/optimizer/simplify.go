package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

type binaryRule func(*ast.BinaryExpr) ast.Expr

var binaryRules = []binaryRule{
	simplifyArithmetic,
	simplifyBoolean,
	canonicalizeBinary,
	reassociateBinary,
}

func simplifyBinary(n *ast.BinaryExpr) ast.Expr {
	for _, rule := range binaryRules {
		if expr := rule(n); expr != nil {
			return expr
		}
	}

	return nil
}
