package sema

import (
	"strings"

	"github.com/azin-lang/Azin/pkg/ast"
	types2 "github.com/azin-lang/Azin/pkg/types"
)

func mangleFunctionName(fn *ast.FuncStmt) string {
	if len(fn.Params) == 0 {
		return fn.Name.Value + "__unit"
	}

	var name strings.Builder
	name.WriteString(fn.Name.Value)
	for _, param := range fn.Params {
		if param.SynType != nil {
			name.WriteString("__" + param.SynType.Value)
		}
	}

	return name.String()
}

func (a *Analyzer) assignFunctionCNames() {
	scope := a.currentScope()

	for name, overloads := range scope.Functions {
		if len(overloads) == 1 {
			overloads[0].Function.CName = name
			continue
		}

		for _, overload := range overloads {
			overload.Function.CName = mangleFunctionName(overload.Function)
		}
	}
}

func (a *Analyzer) resolveCallExpr(n *ast.CallExpr) *Symbol {
	id, ok := n.Callee.(*ast.Identifier)
	if !ok {
		return nil
	}

	overloads := a.lookupFunctions(id.Value)
	if len(overloads) == 0 {
		if !isBuiltin(id.Value) {
			a.errorf(n.Callee, "undefined function: %s", id.Value)
		}
		n.ResolvedName = id.Value
		return nil
	}

	var sameArity bool
	for _, overload := range overloads {
		if len(n.Args) == len(overload.Function.Params) {
			sameArity = true
			break
		}
	}

	if !sameArity {
		a.errorf(n.Callee, "wrong number of arguments to %s", id.Value)
		return nil
	}

	argTypes := make([]*types2.TypeInfo, len(n.Args))
	for i, arg := range n.Args {
		argTypes[i] = a.inferExprType(arg)
	}

	sym := a.resolveCallOverload(id.Value, argTypes)
	if sym == nil {
		a.errorf(n.Callee, "no matching overload for %s", id.Value)
		return nil
	}

	n.ResolvedName = sym.Function.CName
	return sym
}

func (a *Analyzer) verifyResolvedCalls(program *ast.Program) {
	var visitExpr func(ast.Expr)
	var visitStmt func(ast.Stmt)

	visitExpr = func(expr ast.Expr) {
		switch n := expr.(type) {
		case nil, *ast.BadExpr:
			return

		case *ast.CallExpr:
			if n.ResolvedName == "" {
				a.errorf(n.Callee, "internal compiler error: unresolved function call")
			}

			visitExpr(n.Callee)
			for _, arg := range n.Args {
				visitExpr(arg)
			}

		case *ast.BinaryExpr:
			visitExpr(n.Left)
			visitExpr(n.Right)

		case *ast.MemberExpr:
			visitExpr(n.Object)
		}
	}

	visitStmt = func(stmt ast.Stmt) {
		switch n := stmt.(type) {
		case *ast.BadStmt, *ast.ImportCStmt, *ast.ImportStmt, *ast.StructStmt, *ast.EnumStmt, *ast.StopStmt:
			return

		case *ast.FuncStmt:
			for _, stmt := range n.Body {
				visitStmt(stmt)
			}

		case *ast.ReturnStmt:
			visitExpr(n.Value)

		case *ast.VarStmt:
			visitExpr(n.Value)

		case *ast.IfStmt:
			visitExpr(n.Condition)
			for _, stmt := range n.Then {
				visitStmt(stmt)
			}
			for _, stmt := range n.Else {
				visitStmt(stmt)
			}

		case *ast.LoopStmt:
			for _, stmt := range n.Body {
				visitStmt(stmt)
			}

		case *ast.AssignmentStmt:
			visitExpr(n.Left)
			visitExpr(n.Value)

		case *ast.ExpressionStmt:
			visitExpr(n.Expression)

		case *ast.DeferStmt:
			visitExpr(n.Call)
		}
	}

	for _, stmt := range program.Statements {
		visitStmt(stmt)
	}
}

func (a *Analyzer) resolveCallOverload(name string, argTypes []*types2.TypeInfo) *Symbol {
	overloads := a.lookupFunctions(name)
	if len(overloads) == 0 {
		return nil
	}

	var matches []*Symbol

	for _, overload := range overloads {
		if len(argTypes) != len(overload.Function.Params) {
			continue
		}

		match := true
		for i, got := range argTypes {
			want := overload.Function.Params[i].SemaType

			if want == nil || want.IsUnknown() {
				a.errorf(overload.Function.Params[i].SynType, "internal compiler error: parameter type not inferred")
			}

			if !types2.IsAssignable(got, want) {
				match = false
				break
			}
		}

		if match {
			matches = append(matches, overload)
		}
	}

	if len(matches) == 1 {
		return matches[0]
	}

	return nil
}

func isBuiltin(name string) bool {
	switch name {
	case "printf", "fprintf", "sprintf", "snprintf",
		"scanf", "sscanf",
		"malloc", "calloc", "realloc", "free", "exit",
		"strlen", "strcpy", "strcmp", "memset", "memcpy":
		return true
	}
	return false
}
