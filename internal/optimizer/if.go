package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func optimizeIf(n *ast.IfStmt) []ast.Stmt {
	if n.Condition != nil {
		n.Condition = optimizeExpr(n.Condition)
	}

	// Constant folding: if condition is a literal, return the appropriate branch
	if b, ok := n.Condition.(*ast.BooleanLiteral); ok {
		if b.Value {
			return optimizeStatements(n.Then)
		}
		return optimizeStatements(n.Else)
	}

	n.Then = optimizeStatements(n.Then)
	n.Else = optimizeStatements(n.Else)

	// Dead code: empty branches with pure condition
	if len(n.Then) == 0 && len(n.Else) == 0 {
		// Pure condition: delete entirely
		if isPure(n.Condition) {
			return nil
		}

		// Impure condition: keep only the side-effect (convert to expression stmt)
		return []ast.Stmt{
			&ast.ExpressionStmt{Token: n.Token, Expression: n.Condition},
		}
	}

	// If the Then branch guarantees an exit, the Else block can be executed sequentially after the if statement
	if len(n.Else) > 0 && blockIsTerminal(n.Then) {
		// Pre-allocate slice capacity to avoid reallocation
		unnested := make([]ast.Stmt, 0, 1+len(n.Else))

		unnested = append(unnested, n)
		unnested = append(unnested, n.Else...)

		// Remove the else branch from the AST since it's now unnested
		n.Else = nil

		return unnested
	}

	return tryTailMerge(n)
}

func tryTailMerge(n *ast.IfStmt) []ast.Stmt {
	if len(n.Then) == 0 || len(n.Else) == 0 {
		return []ast.Stmt{n}
	}

	lastThen := n.Then[len(n.Then)-1]
	lastElse := n.Else[len(n.Else)-1]

	// If both end in the same return value
	if t, okT := lastThen.(*ast.ReturnStmt); okT {
		if e, okE := lastElse.(*ast.ReturnStmt); okE {
			if exprEqual(t.Value, e.Value) {
				// Drop the returns from the branches
				n.Then = n.Then[:len(n.Then)-1]
				n.Else = n.Else[:len(n.Else)-1]

				// Return the if statement followed by the shared return
				return []ast.Stmt{n, t}
			}
		}
	}
	return []ast.Stmt{n}
}
