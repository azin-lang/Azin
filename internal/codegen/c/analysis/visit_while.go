package analysis

import "github.com/azin-lang/Azin/pkg/ast"

func (a *Analyzer) visitWhile(fn string, stmt *ast.WhileStmt) {
	a.visitExpr(fn, stmt.Condition)

	for _, child := range stmt.Body {
		a.visitStmt(fn, child)
	}
}
