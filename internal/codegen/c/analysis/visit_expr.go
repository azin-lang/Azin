package analysis

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/ast"
)

func (a *Analyzer) visitExpr(fn string, expr ast.Expr) {
	switch n := expr.(type) {

	case nil:
		return

	case *ast.BadExpr:
		return

	case *ast.Identifier:
		a.useVariable(fn, n.Value)

	case *ast.BooleanLiteral:
		a.Transpiler.RequireInclude("stdbool.h")

	case *ast.IntegerLiteral,
		*ast.FloatLiteral,
		*ast.StringLiteral,
		*ast.CharacterLiteral:
		return

	case *ast.BinaryExpr:
		a.visitExpr(fn, n.Left)
		a.visitExpr(fn, n.Right)

	case *ast.CallExpr:
		a.visitCall(fn, n)

	case *ast.MemberExpr:
		a.visitMember(fn, n)

	default:
		panic(fmt.Sprintf(
			"unsupported expression %T",
			expr,
		))
	}
}
