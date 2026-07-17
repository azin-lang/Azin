package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func foldBooleanBoolean(left *ast.BooleanLiteral, op token.Token, right *ast.BooleanLiteral) ast.Expr {
	switch op.Kind {
	case token.LogicalAnd:
		return boolLit(left.Value && right.Value)
	case token.LogicalOr:
		return boolLit(left.Value || right.Value)
	case token.EqualEqual:
		return boolLit(left.Value == right.Value)
	case token.BangEqual:
		return boolLit(left.Value != right.Value)
	default:
		return nil
	}
}
