package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) RemoveUnusedVariables(program *ast.Program) {
	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncStmt); ok {
			a.removeVariables(fn)
		}
	}
}

func (a *Analyzer) removeVariables(fn *ast.FuncStmt) {
	usage := a.Variables[FunctionName(fn)]
	if usage == nil {
		return
	}

	fn.Body = a.pruneBlock(FunctionName(fn), fn.Body, usage)
}

func (a *Analyzer) pruneBlock(fnName string, stmts []ast.Stmt, usage map[string]int) []ast.Stmt {
	out := stmts[:0]
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.VarStmt:
			if usage[s.Name.Value] > 0 || hasSideEffects(s.Value) {
				if s.Value != nil {
					a.visitExpr(fnName, s.Value)
				}
				out = append(out, stmt)

				if s.Type != nil {
					a.MarkTypeUsed(s.Type.Value)
				}
			}
		case *ast.IfStmt:
			s.Then = a.pruneBlock(fnName, s.Then, usage)
			s.Else = a.pruneBlock(fnName, s.Else, usage)

			if len(s.Then) == 0 && len(s.Else) == 0 && !hasSideEffects(s.Condition) {
				continue
			}
			out = append(out, stmt)
		case *ast.LoopStmt:
			s.Body = a.pruneBlock(fnName, s.Body, usage)
			out = append(out, stmt)
		default:
			out = append(out, stmt)
		}
	}
	return out
}

func hasSideEffects(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.CallExpr:
		return true
	case *ast.BinaryExpr:
		return hasSideEffects(e.Left) || hasSideEffects(e.Right)
	case *ast.MemberExpr:
		return hasSideEffects(e.Object)
	default:
		return false
	}
}
