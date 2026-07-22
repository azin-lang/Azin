package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) ComputeReachability() {
	if _, ok := a.Functions["main"]; ok {
		a.reach("main")
	}
}

func (a *Analyzer) BuildCallGraph() {
	for name, info := range a.Functions {
		if info.Stmt != nil {
			a.walkStmtsForCalls(name, info.Stmt.Body)
		}
	}
}

func (a *Analyzer) addCall(from, to string) {
	if from == "" || to == "" {
		return
	}
	if a.Calls[from] == nil {
		a.Calls[from] = make(map[string]struct{})
	}
	a.Calls[from][to] = struct{}{}
}

func (a *Analyzer) walkStmtsForCalls(caller string, stmts []ast.Stmt) {
	for _, stmt := range stmts {
		a.walkStmtForCalls(caller, stmt)
	}
}

func (a *Analyzer) walkStmtForCalls(caller string, stmt ast.Stmt) {
	if stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *ast.ExpressionStmt:
		a.walkExprForCalls(caller, s.Expression)
	case *ast.VarStmt:
		a.walkExprForCalls(caller, s.Value)
	case *ast.ReturnStmt:
		a.walkExprForCalls(caller, s.Value)
	case *ast.AssignmentStmt:
		a.walkExprForCalls(caller, s.Left)
		a.walkExprForCalls(caller, s.Value)
	case *ast.IfStmt:
		a.walkExprForCalls(caller, s.Condition)
		a.walkStmtsForCalls(caller, s.Then)
		a.walkStmtsForCalls(caller, s.Else)
	case *ast.LoopStmt:
		a.walkStmtsForCalls(caller, s.Body)
	}
}

func (a *Analyzer) walkExprForCalls(caller string, expr ast.Expr) {
	walkExpr(expr, func(e ast.Expr) bool {
		if call, ok := e.(*ast.CallExpr); ok {
			name := call.ResolvedName
			if name != "" {
				if _, isBuiltin := LookupBuiltin(name); !isBuiltin {
					if _, isFunc := a.Functions[name]; isFunc {
						a.addCall(caller, name)
					}
				}
			}
		}
		return true // continue walking children
	})
}

func (a *Analyzer) reach(name string) {
	if _, ok := a.ReachableFunctions[name]; ok {
		return
	}
	a.ReachableFunctions[name] = struct{}{}
	for child := range a.Calls[name] {
		a.reach(child)
	}
}
