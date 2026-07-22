package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) RemoveUnusedFunctions(
	program *ast.Program,
) {
	out := program.Statements[:0]

	for _, stmt := range program.Statements {

		fn, ok := stmt.(*ast.FuncStmt)

		if !ok {
			out = append(out, stmt)
			continue
		}

		name := FunctionName(fn)

		if name != "main" {

			if _, ok := a.Reachable[name]; !ok {
				continue
			}
		}

		out = append(out, stmt)
	}

	program.Statements = out
}
