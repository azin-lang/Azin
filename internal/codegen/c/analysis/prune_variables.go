package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) RemoveUnusedVariables(program *ast.Program) {
	for _, stmt := range program.Statements {
		fn, ok := stmt.(*ast.FuncStmt)

		if !ok {
			continue
		}

		a.removeVariables(fn)
	}
}

func (a *Analyzer) removeVariables(fn *ast.FuncStmt) {
	name := FunctionName(fn)

	usage := a.Variables[name]

	if usage == nil {
		return
	}

	out := fn.Body[:0]

	for _, stmt := range fn.Body {
		v, ok := stmt.(*ast.VarStmt)

		if !ok || usage[v.Name.Value] > 0 {
			out = append(out, stmt)
		}
	}

	fn.Body = out
}
