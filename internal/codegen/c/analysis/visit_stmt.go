//nolint:unused
package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitStmt(current string, stmt ast.Stmt) {
	if stmt == nil {
		return
	}
	switch n := stmt.(type) {
	case nil, *ast.BadStmt, *ast.StopStmt:
		return
	case *ast.ImportCStmt:
		a.requireImport(n.Path.Value)
	case *ast.EnumStmt:
		a.visitEnum(n)
	case *ast.FuncStmt:
		a.visitFunction(n)
	case *ast.ReturnStmt:
		a.visitExpr(current, n.Value)
	case *ast.VarStmt:
		a.visitVar(current, n)
	case *ast.IfStmt:
		a.visitIf(current, n)
	case *ast.LoopStmt:
		a.visitLoop(current, n)
	case *ast.AssignmentStmt:
		a.visitAssignment(current, n)
	case *ast.ExpressionStmt:
		a.visitExpr(current, n.Expression)
	}
}

func (a *Analyzer) visitCallGraphStmt(fn string, stmt ast.Stmt) {
	switch n := stmt.(type) {
	case *ast.ExpressionStmt:
		a.visitCallGraphExpr(fn, n.Expression)
	case *ast.VarStmt:
		a.visitCallGraphExpr(fn, n.Value)
	case *ast.ReturnStmt:
		a.visitCallGraphExpr(fn, n.Value)
	case *ast.AssignmentStmt:
		a.visitCallGraphExpr(fn, n.Value)
	case *ast.IfStmt:
		a.visitCallGraphExpr(fn, n.Condition)
		for _, s := range n.Then {
			a.visitCallGraphStmt(fn, s)
		}
		for _, s := range n.Else {
			a.visitCallGraphStmt(fn, s)
		}
	case *ast.LoopStmt:
		for _, s := range n.Body {
			a.visitCallGraphStmt(fn, s)
		}
	}
}
