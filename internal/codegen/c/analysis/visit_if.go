package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitIf(fn string, stmt *ast.IfStmt) {
	a.visitExpr(fn, stmt.Condition)
	for _, child := range stmt.Then {
		a.visitStmt(fn, child)
	}
	for _, child := range stmt.Else {
		a.visitStmt(fn, child)
	}
}
