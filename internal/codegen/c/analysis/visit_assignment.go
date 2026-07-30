package analysis

import (
	"github.com/azin-lang/Azin/pkg/ast"
)

func (a *Analyzer) visitAssignment(fn string, stmt *ast.AssignmentStmt) {
	a.visitExpr(fn, stmt.Left)
	a.visitExpr(fn, stmt.Value)
}
