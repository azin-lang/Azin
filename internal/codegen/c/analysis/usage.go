package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) AnalyzeUsage(program *ast.Program) {
	for _, stmt := range program.Statements {
		fn, ok := stmt.(*ast.FuncStmt)
		if !ok {
			continue
		}
		name := FunctionName(fn)
		if _, reachable := a.ReachableFunctions[name]; !reachable {
			continue
		}
		a.visitFunction(fn)
	}
}
