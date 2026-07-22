package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitVar(fn string, stmt *ast.VarStmt) {
	a.registerVariable(fn, stmt.Name.Value)
}
