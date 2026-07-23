package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func (o *Optimizer) optimizeIf(n *ast.IfStmt) []ast.Stmt {
	if n.Condition != nil {
		n.Condition = o.optimizeExpr(n.Condition)
	}

	if b, ok := n.Condition.(*ast.BooleanLiteral); ok {
		if b.Value {
			o.Enter()
			res := o.optimizeStatements(n.Then)
			o.Leave()
			return res
		}
		o.Enter()
		res := o.optimizeStatements(n.Else)
		o.Leave()
		return res
	}

	o.Enter()
	n.Then = o.optimizeStatements(n.Then)
	o.Leave()

	o.Enter()
	n.Else = o.optimizeStatements(n.Else)
	o.Leave()

	// Try to pull common returns out of the branches.
	var tailReturn *ast.ReturnStmt
	if ret := tryTailMerge(n); ret != nil {
		tailReturn = ret
	}

	// Both branches became empty.
	if len(n.Then) == 0 && len(n.Else) == 0 {
		if isPure(n.Condition) {
			if tailReturn != nil {
				return []ast.Stmt{tailReturn}
			}
			return nil
		}

		stmts := []ast.Stmt{
			&ast.ExpressionStmt{
				Token:      n.Token,
				Expression: n.Condition,
			},
		}

		if tailReturn != nil {
			stmts = append(stmts, tailReturn)
		}

		return stmts
	}

	// Unnest else after terminal then.
	if len(n.Else) > 0 && blockIsTerminal(n.Then) {
		unnested := make([]ast.Stmt, 0, 1+len(n.Else))

		unnested = append(unnested, n)
		unnested = append(unnested, n.Else...)
		n.Else = nil

		if tailReturn != nil {
			unnested = append(unnested, tailReturn)
		}

		return unnested
	}

	if tailReturn != nil {
		return []ast.Stmt{n, tailReturn}
	}

	return []ast.Stmt{n}
}

func tryTailMerge(n *ast.IfStmt) *ast.ReturnStmt {
	if len(n.Then) == 0 || len(n.Else) == 0 {
		return nil
	}

	t, okT := n.Then[len(n.Then)-1].(*ast.ReturnStmt)
	e, okE := n.Else[len(n.Else)-1].(*ast.ReturnStmt)

	if !okT || !okE || !exprEqual(t.Value, e.Value) {
		return nil
	}

	n.Then = n.Then[:len(n.Then)-1]
	n.Else = n.Else[:len(n.Else)-1]

	return t
}
