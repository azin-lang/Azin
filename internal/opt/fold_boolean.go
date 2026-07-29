package opt

import (
	"github.com/azin-lang/Azin/pkg/ast"
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func foldBoolean(left *ast.BooleanLiteral, op token2.Token, right *ast.BooleanLiteral) ast.Expr {
	switch op.Kind {
	case token2.LogicalAnd:
		return boolLit(left.Value && right.Value)
	case token2.LogicalOr:
		return boolLit(left.Value || right.Value)
	case token2.EqualEqual:
		return boolLit(left.Value == right.Value)
	case token2.BangEqual:
		return boolLit(left.Value != right.Value)
	default:
		return nil
	}
}
