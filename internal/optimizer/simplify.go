package optimizer

import "github.com/azin-lang/Azin/internal/ast"

type binaryRule func(*ast.BinaryExpr) ast.Expr

var binaryRules = []binaryRule{
	canonicalizeBinary,
	reassociateBinary,
	simplifyArithmetic,
	simplifyBoolean,
}

func simplifyBinary(n *ast.BinaryExpr) ast.Expr {
	var current ast.Expr = n

	for i := range binaryRules {
		bin, ok := current.(*ast.BinaryExpr)
		if !ok {
			// If a previous rule folded this into a non-binary expression, stop.
			break
		}

		if result := binaryRules[i](bin); result != nil {
			current = result
		}
	}

	if current == n {
		return nil
	}
	return current
}
