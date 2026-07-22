package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitEnum(e *ast.EnumStmt) {
	a.Enums[e.Name.Value] = e
	a.Transpiler.SetEnum(e.Name.Value)
}
