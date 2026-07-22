package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) Walk(program *ast.Program) {
	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncStmt); ok {
			a.visitFunction(fn)
			continue
		}
		a.visitStmt("", stmt)
	}
}

func (a *Analyzer) WalkMain(program *ast.Program) {
	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncStmt); ok && FunctionName(fn) == "main" {
			a.visitFunction(fn)
		} else if !ok {
			a.visitStmt("", stmt)
		}
	}
}
