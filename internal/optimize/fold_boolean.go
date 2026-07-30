package optimize

import (
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/token"
)

func foldBoolean(left *ast.BooleanLiteral, op token.Token, right *ast.BooleanLiteral) ast.Expr {
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
