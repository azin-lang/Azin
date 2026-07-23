package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func (o *Optimizer) optimizeStatements(stmts []ast.Stmt) []ast.Stmt {
	if len(stmts) == 0 {
		return stmts
	}

	var out []ast.Stmt

	for i := range stmts {
		stmt := stmts[i]
		if stmt == nil {
			continue
		}

		optStmts := o.optimizeStatement(stmt)
		if len(optStmts) == 0 {
			continue
		}

		if out == nil {
			out = make([]ast.Stmt, 0, len(stmts))
		}

		for j := range optStmts {
			optStmt := optStmts[j]
			out = append(out, optStmt)

			if isTerminal(optStmt) {
				return o.filterDead(out)
			}
		}
	}
	return o.filterDead(out)
}

func (o *Optimizer) optimizeStatement(stmt ast.Stmt) []ast.Stmt {
	switch n := stmt.(type) {
	case *ast.IfStmt:
		return o.optimizeIf(n)
	case *ast.LoopStmt:
		return o.optimizeLoop(n)
	case *ast.ExpressionStmt:
		return o.optimizeExpressionStmt(n)
	case *ast.FuncStmt:
		o.optimizeFunction(n)
	case *ast.ReturnStmt:
		o.optimizeReturn(n)
	case *ast.VarStmt:
		o.optimizeVariable(n)
	case *ast.AssignmentStmt:
		o.optimizeAssignment(n)
	}
	return []ast.Stmt{stmt}
}

func (o *Optimizer) optimizeLoop(n *ast.LoopStmt) []ast.Stmt {
	o.currentScope.ClearAll()

	o.Enter()
	n.Body = o.optimizeStatements(n.Body)
	o.Leave()

	if len(n.Body) == 0 || !canUnwrapLoop(n.Body) {
		return []ast.Stmt{n}
	}

	last := n.Body[len(n.Body)-1]
	switch last.(type) {
	case *ast.ReturnStmt:
		return n.Body
	case *ast.StopStmt:
		return n.Body[:len(n.Body)-1]
	}
	return []ast.Stmt{n}
}

func (o *Optimizer) optimizeFunction(n *ast.FuncStmt) {
	o.Enter()
	n.Body = o.optimizeStatements(n.Body)
	o.Leave()
}

func (o *Optimizer) optimizeReturn(n *ast.ReturnStmt) {
	if n.Value != nil {
		n.Value = o.optimizeExpr(n.Value)
	}
}

func (o *Optimizer) optimizeVariable(n *ast.VarStmt) {
	if n.Value != nil {
		n.Value = o.optimizeExpr(n.Value)

		o.currentScope.Invalidate(n.Name.Value)
		if isCopyable(n.Value) {
			o.currentScope.SetValue(n.Name.Value, n.Value)
		}

		// DSE: Record this declaration as the most recent store
		o.currentScope.lastStore[n.Name.Value] = n
		o.currentScope.modified[n.Name.Value] = true
	}
}

func (o *Optimizer) optimizeAssignment(n *ast.AssignmentStmt) {
	if _, isId := n.Left.(*ast.Identifier); !isId {
		n.Left = o.optimizeExpr(n.Left)
	}
	n.Value = o.optimizeExpr(n.Value)

	if id, ok := n.Left.(*ast.Identifier); ok {
		// DSE: If there's an active, unused store to this exact variable in this scope, kill it!
		if prev, exists := o.currentScope.lastStore[id.Value]; exists {
			// Do not kill VarStmt. Removing a variable declaration causes
			// the C compiler to fail with an "undeclared variable" error.
			if _, isVarDecl := prev.(*ast.VarStmt); !isVarDecl {
				o.deadStores[prev] = true
			}
		}

		// Break any existing aliases to this variable before re-assigning
		o.currentScope.Invalidate(id.Value)
		if isCopyable(n.Value) {
			o.currentScope.SetValue(id.Value, n.Value)
		}

		// DSE: Record this assignment as the most recent store
		o.currentScope.lastStore[id.Value] = n
		o.currentScope.modified[id.Value] = true
	}
}

func (o *Optimizer) optimizeExpressionStmt(n *ast.ExpressionStmt) []ast.Stmt {
	if n.Expression == nil {
		return nil
	}
	n.Expression = o.optimizeExpr(n.Expression)
	if isPure(n.Expression) {
		return nil
	}
	return []ast.Stmt{n}
}

// isTerminal checks if a single statement halts execution flow.
func isTerminal(stmt ast.Stmt) bool {
	switch s := stmt.(type) {
	case *ast.ReturnStmt, *ast.StopStmt:
		return true
	case *ast.IfStmt:
		// If both branches halt, the if statement itself guarantees a halt.
		return blockIsTerminal(s.Then) && blockIsTerminal(s.Else)
	}
	return false
}

// blockIsTerminal checks if a slice of statements is guaranteed to halt.
// Because optimizeStatements strips dead code, a terminal statement
// will always be exactly at the last index.
func blockIsTerminal(stmts []ast.Stmt) bool {
	if len(stmts) == 0 {
		return false
	}
	return isTerminal(stmts[len(stmts)-1])
}

func canUnwrapLoop(body []ast.Stmt) bool {
	if len(body) == 0 {
		return false
	}

	for i := 0; i < len(body)-1; i++ {
		if !isSimpleStmt(body[i]) {
			return false
		}
	}

	return isTerminal(body[len(body)-1])
}

func isSimpleStmt(stmt ast.Stmt) bool {
	switch stmt.(type) {
	case *ast.ExpressionStmt,
		*ast.VarStmt,
		*ast.AssignmentStmt:
		return true
	default:
		return false
	}
}

func (o *Optimizer) filterDead(stmts []ast.Stmt) []ast.Stmt {
	if len(o.deadStores) == 0 {
		return stmts
	}

	var filtered []ast.Stmt
	for _, s := range stmts {
		if !o.deadStores[s] {
			filtered = append(filtered, s)
			continue
		}

		// The store is dead, but we must preserve side effects (like function calls).
		var rhs ast.Expr
		switch n := s.(type) {
		case *ast.VarStmt:
			rhs = n.Value
		case *ast.AssignmentStmt:
			rhs = n.Value
		}

		if rhs != nil && !isPure(rhs) {
			filtered = append(filtered, &ast.ExpressionStmt{
				Expression: rhs,
			})
		}
	}
	return filtered
}
