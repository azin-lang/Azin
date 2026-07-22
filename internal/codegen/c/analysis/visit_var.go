package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitVar(
	fn string,
	stmt *ast.VarStmt,
) {
	if stmt.Type != nil {
		a.markType(stmt.Type.Value)
	}

	a.registerVariable(
		fn,
		stmt.Name.Value,
	)

	if stmt.Value != nil {
		a.visitExpr(fn, stmt.Value)
	}
}
