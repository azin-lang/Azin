//nolint:goconst
package sema

import (
	"slices"

	"github.com/azin-lang/Azin/pkg/ast"
	types2 "github.com/azin-lang/Azin/pkg/types"
)

func (a *Analyzer) visitStatement(stmt ast.Stmt) {
	switch n := stmt.(type) {

	case *ast.BadStmt:
		return

	case *ast.FuncStmt:
		old := a.currentFunction
		a.currentFunction = n
		defer func() {
			a.currentFunction = old
		}()

		a.pushScope()

		// Register parameters.
		for _, param := range n.Params {
			if a.checkEnumShadow(param.Name) {
				continue
			}

			if !param.SemaType.IsKnown() {
				a.errorf(param.SynType, "internal compiler error: parameter type is null or not inferred")
				param.SemaType = types2.ErrorType()
			}

			a.declare(&Symbol{
				Name:     param.Name.Value,
				Type:     param.SemaType,
				Kind:     SymbolVariable,
				DeclNode: param.Name,
			})
		}

		if !n.SemaReturnType.IsKnown() {
			a.inferFunctionReturnType(n)

			if sym := a.lookupFunctionSymbol(n); sym != nil {
				sym.Type = n.SemaReturnType
			}
		}

		for _, stmt := range n.Body {
			a.visitStatement(stmt)
		}

		if n.SemaReturnType != nil && !n.SemaReturnType.IsUnit() {
			if !a.blockAllPathsReturn(n.Body) {
				a.errorf(
					n.Name,
					"function '%s' must return %s on all paths",
					n.Name.Value,
					n.SynReturnType.Value,
				)
			}
		}

		a.popScope()

	case *ast.ReturnStmt:
		if a.currentFunction == nil {
			return
		}

		actual := types2.UnitType()
		if n.Value != nil {
			actual = a.inferExprType(n.Value)
		}

		expected := a.currentFunction.SemaReturnType

		if expected != nil && actual != nil && !types2.IsAssignable(actual, expected) {
			var posErr ast.Node
			if n.Value == nil {
				posErr = n
			} else {
				posErr = n.Value
			}
			a.errorf(
				posErr,
				"return type mismatch: expected %s, got %s",
				expected.Name,
				actual.Name,
			)
		}

	case *ast.DeferStmt:
		a.analyzeExpr(n.Call)

	case *ast.VarStmt:
		if n.Value == nil && n.SynType == nil {
			a.errorf(n.Name, "variable '%s' must have a type or an initializer", n.Name.Value)
			return
		}

		varType := types2.UnknownType()
		if n.SynType != nil {
			foundType, err := a.lookupType(n.SynType.Value)
			if err != nil {
				a.errorf(n.SynType, "unknown type: %s", n.SynType.Value)
			}

			varType = foundType
		}

		if n.Value != nil {
			valueType := a.inferExprType(n.Value)
			if varType.IsUnknown() {
				varType = valueType
			} else if !types2.IsAssignable(valueType, varType) {
				a.errorf(
					n.Value,
					"cannot initialize variable '%s' of type '%s' with value of type '%s'",
					n.Name.Value,
					varType.Name,
					valueType.Name,
				)
			}
		}

		n.SemaType = varType

		if a.checkEnumShadow(n.Name) {
			return
		}

		a.declare(&Symbol{
			Name:     n.Name.Value,
			Type:     varType,
			Kind:     SymbolVariable,
			Mutable:  n.Mutable,
			DeclNode: n.Name,
		})

	case *ast.IfStmt:
		cond := a.inferExprType(n.Condition)
		if cond != nil && !cond.IsBool() {
			a.errorf(n.Condition, "if condition must be bool, got %s", cond.Name)
		}

		a.pushScope()

		for _, stmt := range n.Then {
			a.visitStatement(stmt)
		}

		a.popScope()

		a.pushScope()

		for _, stmt := range n.Else {
			a.visitStatement(stmt)
		}

		a.popScope()

	case *ast.WhileStmt:
		a.loopDepth++
		defer func() { a.loopDepth-- }()

		a.pushScope()
		defer a.popScope()

		cond := a.inferExprType(n.Condition)
		if !types2.IsAssignable(cond, types2.BoolType()) {
			a.errorf(n.Condition, "while condition must be bool, got %s", cond.Name)
		}

		for _, stmt := range n.Body {
			a.visitStatement(stmt)
		}

	case *ast.LoopStmt:
		a.loopDepth++
		defer func() { a.loopDepth-- }()

		a.pushScope()
		defer a.popScope()

		for _, stmt := range n.Body {
			a.visitStatement(stmt)
		}

	case *ast.StopStmt:
		if a.loopDepth == 0 {
			a.errorf(n, "'stop' can only be used inside a loop")
		}

	case *ast.AssignmentStmt:
		switch left := n.Left.(type) {

		case *ast.BadExpr:
			// Skip type checking for malformed assignment targets
			return

		case *ast.Identifier:
			sym := a.lookup(left.Value)
			if sym == nil {
				a.errorf(n.Value, "unknown variable: %s", left.Value)
				return
			}

			if sym.Kind != SymbolVariable {
				a.errorf(n.Value, "%s is not a variable", left.Value)
				return
			}

			if !sym.Mutable {
				a.errorf(n.Value, "cannot assign to immutable variable '%s'", left.Value)
				return
			}

			got := a.inferExprType(n.Value)

			if got != nil && sym.Type != nil && !types2.IsAssignable(got, sym.Type) {
				a.errorf(
					n.Value,
					"cannot assign %s to variable '%s' of type %s",
					got.Name,
					left.Value,
					sym.Type.Name,
				)
			}

		case *ast.MemberExpr:
			objectType := a.inferExprType(left.Object)
			if objectType == nil {
				a.errorf(n.Value, "cannot determine type of member access")
				return
			}

			if objectType.Kind != types2.Nominal {
				a.errorf(n.Value, "'%s' is not a struct", objectType.Name)
				return
			}

			if sym := a.lookup(objectType.Name); sym != nil && sym.Kind == SymbolEnum {
				a.errorf(n.Value, "'%s' is an enum and cannot be assigned to", objectType.Name)
				return
			}

			strct := a.lookupStruct(objectType.Name)
			if strct == nil {
				a.errorf(n.Value, "'%s' is not a struct", objectType.Name)
				return
			}

			field := a.lookupField(strct, left.Property.Value)
			if field == nil {
				a.errorf(n.Value, "struct '%s' has no field '%s'", strct.Name.Value, left.Property.Value)
				return
			}

			if !field.Mutable {
				a.errorf(n.Value, "cannot assign to immutable field '%s'", field.Name.Value)
				return
			}

			got := a.inferExprType(n.Value)

			if got != nil && !types2.IsAssignable(got, field.SemaType) {
				a.errorf(
					n.Value, "cannot assign %s to field '%s' of type %s",
					got.Name,
					field.Name.Value,
					field.SemaType.Name,
				)
			}

		default:
			a.errorf(n.Value, "left side of assignment is not assignable")
		}

	case *ast.ExpressionStmt:
		a.analyzeExpr(n.Expression)
	}
}

func (a *Analyzer) findReturnExprType(stmt ast.Stmt) *types2.TypeInfo {
	switch n := stmt.(type) {
	case *ast.ReturnStmt:
		if n.Value == nil {
			return types2.UnitType()
		}
		return a.inferExprType(n.Value)

	case *ast.IfStmt:
		for _, s := range n.Then {
			if t := a.findReturnExprType(s); t != nil {
				return t
			}
		}
		for _, s := range n.Else {
			if t := a.findReturnExprType(s); t != nil {
				return t
			}
		}

	case *ast.LoopStmt:
		for _, s := range n.Body {
			if t := a.findReturnExprType(s); t != nil {
				return t
			}
		}

	case *ast.DeferStmt:
	}

	return nil
}

func (a *Analyzer) blockAllPathsReturn(stmts []ast.Stmt) bool {
	return slices.ContainsFunc(stmts, a.stmtAllPathsReturn)
}

func (a *Analyzer) stmtAllPathsReturn(stmt ast.Stmt) bool {
	switch n := stmt.(type) {
	case *ast.ReturnStmt:
		return true

	case *ast.StopStmt:
		return true

	case *ast.IfStmt:
		if len(n.Else) == 0 {
			return false
		}
		return a.blockAllPathsReturn(n.Then) && a.blockAllPathsReturn(n.Else)

	case *ast.LoopStmt:
		return a.blockAllPathsReturn(n.Body)

	case *ast.BadStmt:
		return false

	case *ast.DeferStmt:
		return false

	default:
		return false
	}
}

func (a *Analyzer) inferFunctionReturnType(fn *ast.FuncStmt) {
	if fn.SemaReturnType.IsKnown() {
		// The return type is already known, so we don't need to infer it.
		return
	}

	sym := a.lookupFunctionSymbol(fn)
	if sym != nil {
		if sym.Inferring {
			return
		}
		sym.Inferring = true
		defer func() {
			sym.Inferring = false
		}()
	}

	for _, stmt := range fn.Body {
		if typ := a.findReturnExprType(stmt); typ != nil {
			fn.SemaReturnType = typ

			if sym != nil {
				sym.Type = typ
			}

			return
		}
	}

	fn.SemaReturnType = types2.UnitType()

	if sym != nil {
		sym.Type = fn.SemaReturnType
	}
}
