package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitLoop(fn string, stmt *ast.LoopStmt) {
	for _, child := range stmt.Body {
		a.visitStmt(fn, child)
	}
}
