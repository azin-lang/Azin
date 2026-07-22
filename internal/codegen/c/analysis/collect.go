package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) CollectFunctions(program *ast.Program) {
	for index, stmt := range program.Statements {
		fn, ok := stmt.(*ast.FuncStmt)
		if !ok {
			continue
		}

		name := FunctionName(fn)

		a.Functions[name] = FuncInfo{
			Stmt:  fn,
			Index: index,
		}

		a.Transpiler.SetFunctionIndex(name, index)
	}
}

func FunctionName(fn *ast.FuncStmt) string {
	if fn.CName != "" {
		return fn.CName
	}

	return fn.Name.Value
}
