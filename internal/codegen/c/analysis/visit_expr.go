package analysis

import (

	"github.com/azin-lang/Azin/internal/ast"
)

func (a *Analyzer) visitExpr(fn string, expr ast.Expr) {
	switch n := expr.(type) {
	case nil, *ast.BadExpr, *ast.IntegerLiteral, *ast.FloatLiteral, *ast.StringLiteral, *ast.CharacterLiteral:
		return
	case *ast.Identifier:
		a.useVariable(fn, n.Value)
	case *ast.BooleanLiteral:
		a.Transpiler.RequireInclude("stdbool.h")
	case *ast.BinaryExpr:
		a.visitExpr(fn, n.Left)
		a.visitExpr(fn, n.Right)
	case *ast.CallExpr:
		a.visitCall(fn, n)
	case *ast.MemberExpr:
		a.visitMember(fn, n)
	default:
		a.ReportError("unsupported expression %T", expr)
	}
}

func (a *Analyzer) visitCallGraphExpr(current string, expr ast.Expr) {
	switch e := expr.(type) {
	case nil:
		return
	case *ast.CallExpr:
		name := e.ResolvedName
		if _, builtin := LookupBuiltin(name); builtin {
			return
		}
		if _, ok := a.Functions[name]; ok {
			a.addCall(current, name)
		}
		for _, arg := range e.Args {
			a.visitCallGraphExpr(current, arg)
		}
	case *ast.BinaryExpr:
		a.visitCallGraphExpr(current, e.Left)
		a.visitCallGraphExpr(current, e.Right)
	}
}
