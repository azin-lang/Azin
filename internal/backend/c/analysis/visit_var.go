package analysis

import (
	"github.com/azin-lang/Azin/pkg/ast"
)

func (a *Analyzer) visitVar(fn string, stmt *ast.VarStmt) {
	a.registerVariable(fn, stmt.Name.Value)
	a.visitExpr(fn, stmt.Value)
}
