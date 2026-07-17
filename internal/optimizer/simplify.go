package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

type binaryRule func(*ast.BinaryExpr) ast.Expr

var binaryRules = []binaryRule{
	canonicalizeBinary,
	reassociateBinary,
	simplifyArithmetic,
	simplifyBoolean,
}

func simplifyBinary(n *ast.BinaryExpr) ast.Expr {
	current := ast.Expr(n)

	// Iterate through rules, allowing them to transform the expression
	// and passing the `result` to the next rule in the pipeline.
	for _, rule := range binaryRules {
		// We need to cast back to *ast.BinaryExpr if the rule returns a new one
		if bin, ok := current.(*ast.BinaryExpr); ok {
			if result := rule(bin); result != nil {
				current = result
			}
		}
	}

	// If 'current' is still the original node and no rule changed it, return nil
	if current == n {
		return nil
	}

	return current
}
