package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) RemoveUnusedFunctions(program *ast.Program) {
	out := program.Statements[:0]
	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncStmt); ok {
			name := FunctionName(fn)
			if name != "main" {
				if _, ok := a.ReachableFunctions[name]; !ok {
					continue // Prune dead function
				}
			}
		}
		out = append(out, stmt)
	}
	program.Statements = out
}
