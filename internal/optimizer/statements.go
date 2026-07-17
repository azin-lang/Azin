package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func optimizeStatements(stmts []ast.Stmt) []ast.Stmt {
	if len(stmts) == 0 {
		return stmts
	}

	var out []ast.Stmt // Delay allocation until required

	for i := range stmts {
		stmt := stmts[i]
		if stmt == nil {
			continue
		}

		optStmts := optimizeStatement(stmt)
		if len(optStmts) == 0 {
			continue
		}

		if out == nil {
			out = make([]ast.Stmt, 0, len(stmts))
		}

		for j := range optStmts {
			optStmt := optStmts[j]
			out = append(out, optStmt)

			switch optStmt.(type) {
			case *ast.ReturnStmt, *ast.StopStmt:
				// Terminal statement reached; abandon dead code
				return out
			}
		}
	}

	return out
}

func optimizeStatement(stmt ast.Stmt) []ast.Stmt {
	switch n := stmt.(type) {
	case *ast.IfStmt:
		return optimizeIf(n)

	case *ast.LoopStmt:
		return optimizeLoop(n)

	case *ast.ExpressionStmt:
		return optimizeExpressionStmt(n)

	case *ast.FuncStmt:
		optimizeFunction(n)

	case *ast.ReturnStmt:
		optimizeReturn(n)

	case *ast.VarStmt:
		optimizeVariable(n)

	case *ast.AssignmentStmt:
		optimizeAssignment(n)

	}

	return []ast.Stmt{stmt}
}

func optimizeLoop(n *ast.LoopStmt) []ast.Stmt {
	n.Body = optimizeStatements(n.Body)

	// TODO: revisit this when we add conditions to loops, so we can discard or simplify loops such as while(false)
	return []ast.Stmt{n}
}

func optimizeFunction(n *ast.FuncStmt) {
	n.Body = optimizeStatements(n.Body)
}

func optimizeReturn(n *ast.ReturnStmt) {
	if n.Value != nil {
		n.Value = optimizeExpr(n.Value)
	}
}

func optimizeVariable(n *ast.VarStmt) {
	if n.Value != nil {
		n.Value = optimizeExpr(n.Value)
	}
}

func optimizeAssignment(n *ast.AssignmentStmt) {
	n.Left = optimizeExpr(n.Left)
	n.Value = optimizeExpr(n.Value)
}

func optimizeExpressionStmt(n *ast.ExpressionStmt) []ast.Stmt {
	if n.Expression != nil {
		n.Expression = optimizeExpr(n.Expression)

		// If the statement has no side effects, the statement is useless.
		// e.g. "5 + 5;" or "x;".
		if isPure(n.Expression) {
			return []ast.Stmt{}
		}
	}

	return []ast.Stmt{n}
}
