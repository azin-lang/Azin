package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

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

	// Invert condition if 'Then' is empty but 'Else' is not.
	// This changes `if (a == b) {} else { foo() }` to `if (a != b) { foo() }`
	if len(n.Then) == 0 && len(n.Else) > 0 {
		if invCond := invertCondition(n.Condition); invCond != nil {
			n.Condition = invCond
			n.Then = n.Else
			n.Else = nil
		}
	}

	// if (a) { if (b) { foo() } } -> if (a && b) { foo() }
	if len(n.Else) == 0 && len(n.Then) == 1 {
		if inner, ok := n.Then[0].(*ast.IfStmt); ok && len(inner.Else) == 0 {
			combined := &ast.BinaryExpr{
				Left:     n.Condition,
				Operator: token.Token{Kind: token.LogicalAnd},
				Right:    inner.Condition,
			}
			// Re-optimize the new condition in case `a && b` can be folded
			n.Condition = o.optimizeExpr(combined)
			n.Then = inner.Then
		}
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

// invertCondition attempts to mathematically negate a condition.
// Returns nil if the condition cannot be safely inverted.
func invertCondition(expr ast.Expr) ast.Expr {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	var inverse token.Kind
	switch bin.Operator.Kind {
	case token.EqualEqual:
		inverse = token.BangEqual
	case token.BangEqual:
		inverse = token.EqualEqual
	case token.Less:
		// Unsafe for floats due to NaN behavior
		if !isNotFloat(bin.Left) {
			return nil
		}
		inverse = token.GreaterEqual
	case token.LessEqual:
		if !isNotFloat(bin.Left) {
			return nil
		}
		inverse = token.Greater
	case token.Greater:
		if !isNotFloat(bin.Left) {
			return nil
		}
		inverse = token.LessEqual
	case token.GreaterEqual:
		if !isNotFloat(bin.Left) {
			return nil
		}
		inverse = token.Less
	default:
		return nil
	}

	return &ast.BinaryExpr{
		Left:     bin.Left,
		Operator: token.Token{Kind: inverse},
		Right:    bin.Right,
	}
}
