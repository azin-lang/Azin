//nolint:goconst
package sema

import (
	"fmt"
	"slices"

	"strings"

	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/diagnostics"
	"github.com/azin-lang/Azin/pkg/token"
	types2 "github.com/azin-lang/Azin/pkg/types"
)

/*
- The Analyzer struct is responsible for performing sema analysis on the AST of a program.
- It maintains a stack of scopes, keeps track of the current function being analyzed,
- and uses a diagnostics engine to report errors and warnings.
*/
type Analyzer struct {
	scopes []*Scope

	currentFunction *ast.FuncStmt
	diag            *diagnostics.Engine

	loopDepth int
}

// New creates a new Analyzer instance with the provided diagnostics engine.
func New(diag *diagnostics.Engine) *Analyzer {
	return &Analyzer{
		diag: diag,
	}
}

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

func (a *Analyzer) analyzeExpr(expr ast.Expr) {
	switch n := expr.(type) {
	case nil, *ast.BadExpr:
		return

	case *ast.Identifier:
		a.lookup(n.Value)

	case *ast.CallExpr:
		a.analyzeExpr(n.Callee)
		for _, arg := range n.Args {
			a.analyzeExpr(arg)
		}
		a.resolveCallExpr(n)

	case *ast.BinaryExpr:
		a.analyzeExpr(n.Left)
		a.analyzeExpr(n.Right)

	case *ast.MemberExpr:
		a.analyzeExpr(n.Object)
	}
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
		case *ast.BadStmt, *ast.ImportCStmt, *ast.StructStmt, *ast.EnumStmt, *ast.StopStmt:
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

			if got == nil || want == nil || !types2.IsAssignable(got, want) {
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

func (a *Analyzer) errorf(node ast.Node, format string, args ...any) {
	pos, length := sourceSpan(node)
	a.diag.ReportError(
		pos,
		int(length),
		format,
		args...,
	)
}

func (a *Analyzer) warningf(node ast.Node, format string, args ...any) {
	pos, length := sourceSpan(node)
	a.diag.ReportWarning(
		pos,
		int(length),
		format,
		args...,
	)
}

func sourceSpan(n ast.Node) (pos token.Position, length uint32) {
	pos = n.Pos()
	switch node := n.(type) {
	case *ast.Identifier:
		return pos, node.Token.Length

	case *ast.IntegerLiteral:
		return pos, node.Token.Length
	case *ast.FloatLiteral:
		return pos, node.Token.Length
	case *ast.StringLiteral:
		return pos, node.Token.Length
	case *ast.CharacterLiteral:
		return pos, node.Token.Length
	case *ast.BooleanLiteral:
		return pos, node.Token.Length
	case *ast.BinaryExpr:
		leftEnd := node.Left.Pos().Offset + spanLen(node.Left)
		rightEnd := node.Right.Pos().Offset + spanLen(node.Right)
		if rightEnd > leftEnd {
			return pos, rightEnd - pos.Offset
		}
		return pos, leftEnd - pos.Offset
	case *ast.MemberExpr:
		objEnd := node.Object.Pos().Offset + spanLen(node.Object)
		propEnd := node.Property.Token.Position.Offset + node.Property.Token.Length
		if objEnd > propEnd {
			return pos, objEnd - pos.Offset
		}
		return pos, propEnd - pos.Offset
	case *ast.CallExpr:
		end := node.Callee.Pos().Offset + spanLen(node.Callee)
		for _, arg := range node.Args {
			argEnd := arg.Pos().Offset + spanLen(arg)
			if argEnd > end {
				end = argEnd
			}
		}
		return pos, end - pos.Offset
	case *ast.VarStmt:
		end := node.Name.Token.Position.Offset + node.Name.Token.Length
		if node.SynType != nil {
			tEnd := node.SynType.Token.Position.Offset + node.SynType.Token.Length
			if tEnd > end {
				end = tEnd
			}
		}
		if node.Value != nil {
			vEnd := node.Value.Pos().Offset + spanLen(node.Value)
			if vEnd > end {
				end = vEnd
			}
		}
		return pos, end - pos.Offset
	case *ast.AssignmentStmt:
		end := node.Left.Pos().Offset + spanLen(node.Left)
		vEnd := node.Value.Pos().Offset + spanLen(node.Value)
		if vEnd > end {
			end = vEnd
		}
		return pos, end - pos.Offset
	}

	l := len(n.TokenLiteral())
	if l > 0 {
		length = uint32(l) //nolint:gosec
	}
	return
}

func spanLen(n ast.Node) uint32 {
	_, length := sourceSpan(n)
	return length
}

// Analyze performs sema analysis on the given AST program.
func (a *Analyzer) Analyze(program *ast.Program) error {
	a.pushScope()
	defer a.popScope()

	a.registerTopLevelSymbols(program)
	a.assignFunctionCNames()

	// Infer function return types before sema analysis.
	// for _, stmt := range program.Statements {
	//	if fn, ok := stmt.(*ast.FuncStmt); ok {
	//		a.inferFunctionReturnType(fn)
	//
	//		sym := a.lookup(fn.Name.Value)
	//		if sym != nil {
	//			sym.Type = fn.ReturnType
	//		}
	//	}
	//}

	for _, stmt := range program.Statements {
		a.visitStatement(stmt)
	}

	a.verifyResolvedCalls(program)

	return a.diag.Err()
}

func (a *Analyzer) registerTopLevelSymbols(program *ast.Program) {
	// Register top-level types first, so that they can be used in function signatures and variable declarations.
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {

		case *ast.StructStmt:
			a.declare(&Symbol{
				Name:   n.Name.Value,
				Kind:   SymbolStruct,
				Struct: n,
			})

			n.SemaType = types2.NominalType(n.Name.Value)

		case *ast.EnumStmt:
			a.declare(&Symbol{
				Name: n.Name.Value,
				Kind: SymbolEnum,
				Enum: n,
			})

			n.SemaType = types2.NominalType(n.Name.Value)

		}
	}

	// Register top-level functions their parameters and return type resolved + Resolving field's types.
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {
		case *ast.StructStmt:
			// Resolve the types of the fields in the struct.
			for _, field := range n.Fields {
				fieldType := a.lookupType(field.SynType.Value)
				if fieldType == nil {
					a.errorf(field.SynType, "unknown type: %s", field.SynType.Value)
					fieldType = types2.ErrorType()
				}

				// Store the TypeInfo in the field node for later phase.
				field.SemaType = fieldType
			}

		case *ast.FuncStmt:
			returnType := types2.UnknownType()
			if n.SynReturnType != nil {
				returnType = a.lookupType(n.SynReturnType.Value)
				if returnType == nil {
					a.errorf(n.SynReturnType, "unknown return type: %s", n.SynReturnType.Value)
					returnType = types2.ErrorType()
				}
			}

			a.declareFunction(&Symbol{
				Name:     n.Name.Value,
				Type:     returnType,
				Kind:     SymbolFunction,
				Function: n,
			})

			for _, param := range n.Params {
				if param.SynType == nil {
					continue
				}
				paramType := a.lookupType(param.SynType.Value)
				if paramType == nil {
					a.errorf(param.SynType, "unknown parameter type: %s", param.SynType.Value)
					paramType = types2.ErrorType()
				}

				// Store the TypeInfo in the parameter node for later phase.
				param.SemaType = paramType
			}

			n.SemaReturnType = returnType
		}
	}
}

func (a *Analyzer) lookupType(name string) *types2.TypeInfo {
	if primitive := types2.GetPrimitiveTypeByName(name); primitive != nil {
		return primitive
	}

	sym := a.lookup(name)
	if sym == nil {
		return nil
	}

	switch sym.Kind {
	case SymbolStruct, SymbolEnum:
		return types2.NominalType(sym.Name)
	default:
		if sym.Type != nil {
			return sym.Type
		}

		fmt.Println("internal compiler error: symbol has no type info:", sym.Name)
		return nil
	}
}

func (a *Analyzer) lookupStruct(name string) *ast.StructStmt {
	sym := a.lookup(name)
	if sym == nil || sym.Kind != SymbolStruct {
		return nil
	}

	return sym.Struct
}

func (a *Analyzer) lookupField(strct *ast.StructStmt, name string) *ast.FieldDecl {
	for _, field := range strct.Fields {
		if field.Name.Value == name {
			return field
		}
	}

	return nil
}

func (a *Analyzer) checkEnumShadow(name *ast.Identifier) bool {
	if sym := a.lookup(name.Value); sym != nil && sym.Kind == SymbolEnum {
		a.errorf(name, "cannot shadow enum '%s' with a variable", name.Value)
		return true
	}
	return false
}

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

			if param.SemaType == nil || param.SemaType.IsUnknown() {
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

		if n.SemaReturnType == nil || n.SemaReturnType.IsUnknown() {
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
			varType = a.lookupType(n.SynType.Value)
			if varType == nil {
				a.errorf(n.SynType, "unknown type: %s", n.SynType.Value)
				varType = types2.ErrorType()
			}
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
	if fn.SemaReturnType != nil && !fn.SemaReturnType.IsUnknown() {
		// The return type has already been inferred or explicitly specified, so we don't need to infer it again.
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
			fn.SynReturnType = &ast.Identifier{Value: typ.Name}
			fn.SemaReturnType = typ

			if sym != nil {
				sym.Type = typ
			}

			return
		}
	}

	fn.SynReturnType = &ast.Identifier{Value: "unit"}
	fn.SemaReturnType = types2.UnitType()

	if sym != nil {
		sym.Type = fn.SemaReturnType
	}
}

func (a *Analyzer) inferExprType(expr ast.Expr) *types2.TypeInfo {
	if expr != nil && expr.Type().IsComplete() {
		// Type has already been inferred, so we can return it directly.
		return expr.Type()
	}

	switch n := expr.(type) {
	case *ast.Identifier:
		resultType := a.lookupType(n.Value)
		if resultType == nil {
			resultType = types2.ErrorType()
			a.errorf(n, "unknown identifier: %s", n.Value)
		}

		n.SemaType = resultType
		return resultType

	case *ast.CallExpr:
		sym := a.resolveCallExpr(n)

		n.SemaReturnType = types2.UnknownType()
		if sym == nil {
			// FIXME: any function returning a C call's result will fail. When are we creating a signature table for the headers?
			// Or better yet, add actual header parsing...
			n.SemaReturnType = types2.ErrorType()
			return types2.ErrorType()
		}

		if sym.Inferring {
			return types2.UnknownType()
		}

		if !sym.Type.IsComplete() {
			// If the function's return type is not yet inferred, we need to infer it now.
			a.inferFunctionReturnType(sym.Function)
		}

		n.SemaReturnType = sym.Type
		return sym.Type

	case *ast.BinaryExpr:
		left := a.inferExprType(n.Left)
		right := a.inferExprType(n.Right)

		if !left.IsValid() || !right.IsValid() {
			n.SemaResultType = types2.ErrorType()
			return types2.ErrorType()
		}

		//nolint:exhaustive
		switch n.Operator.Kind {
		case token.Plus, token.Minus, token.Star, token.Slash:
			if !left.IsNumeric() || !right.IsNumeric() {
				a.errorf(
					n,
					"operator '%s' requires numeric operands",
					n.TokenLiteral(),
				)

				n.SemaResultType = types2.ErrorType()
				return types2.ErrorType()
			}

			if left.IsFloat() || right.IsFloat() {
				n.SemaResultType = types2.FloatType()
				return types2.FloatType()
			}

			n.SemaResultType = types2.IntType()
			return types2.IntType()

		case token.EqualEqual, token.BangEqual,
			token.Less, token.LessEqual,
			token.Greater, token.GreaterEqual:

			if !left.Equals(right) {
				a.errorf(
					n,
					"cannot compare %s with %s",
					left.Name,
					right.Name,
				)
				n.SemaResultType = types2.ErrorType()
				return types2.ErrorType()
			}

			n.SemaResultType = types2.BoolType()
			return types2.BoolType()
		}

		n.SemaResultType = types2.ErrorType()
		return types2.ErrorType()

	case *ast.MemberExpr:
		// the object is a type name used as a namespace, so it must be resolved before type inference
		if id, ok := n.Object.(*ast.Identifier); ok {
			if sym := a.lookup(id.Value); sym != nil && sym.Kind == SymbolEnum {
				for _, variant := range sym.Enum.Variants {
					if variant.Value == n.Property.Value {
						n.SemaResultType = types2.NominalType(sym.Enum.Name.Value)
						return n.SemaResultType
					}
				}

				a.errorf(n.Property, "enum '%s' has no variant '%s'", sym.Enum.Name.Value, n.Property.Value)
				n.SemaResultType = types2.ErrorType()
				return n.SemaResultType
			}
		}

		objectType := a.inferExprType(n.Object)
		if !objectType.IsValid() {
			n.SemaResultType = types2.ErrorType()
			return n.SemaResultType
		}

		if objectType.Kind != types2.Nominal {
			a.errorf(n.Object, "'%s' is not a struct", objectType.Name)
			n.SemaResultType = types2.ErrorType()
			return n.SemaResultType
		}

		strct := a.lookupStruct(objectType.Name)
		if strct == nil {
			a.errorf(n.Object, "'%s' is not a struct", objectType.Name)
			n.SemaResultType = types2.ErrorType()
			return n.SemaResultType
		}

		for _, field := range strct.Fields {
			if field.Name.Value == n.Property.Value {
				n.SemaResultType = field.SemaType
				return n.SemaResultType
			}
		}

		a.errorf(n.Property, "struct '%s' has no field '%s'", strct.Name.Value, n.Property.Value)
		n.SemaResultType = types2.ErrorType()
		return n.SemaResultType
	}

	a.errorf(expr, "internal compiler error: cannot infer type for expression")
	return types2.ErrorType()
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
