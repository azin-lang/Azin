package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitFunction(
	fn *ast.FuncStmt,
) {
	name := FunctionName(fn)

	if fn.ReturnType != nil {
		a.markType(fn.ReturnType.Value)
	}

	for _, param := range fn.Params {
		if param.Type != nil {
			a.markType(param.Type.Value)
		}
	}

	for _, stmt := range fn.Body {
		a.visitStmt(name, stmt)
	}
}
