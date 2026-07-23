package sema

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/token"
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

	name := fn.Name.Value
	for _, param := range fn.Params {
		name += "__" + param.Type.Value
	}

	return name
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

	argTypes := make([]*ast.Identifier, len(n.Args))
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
		}
	}

	for _, stmt := range program.Statements {
		visitStmt(stmt)
	}
}

func (a *Analyzer) resolveCallOverload(name string, argTypes []*ast.Identifier) *Symbol {
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
			want := overload.Function.Params[i].Type

			if got == nil || want == nil || got.Value != want.Value {
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
	a.diag.ReportError(
		node.Pos(),
		uint32(len(node.TokenLiteral())),
		format,
		args...,
	)
}

func (a *Analyzer) warningf(node ast.Node, format string, args ...any) {
	a.diag.ReportWarning(
		node.Pos(),
		uint32(len(node.TokenLiteral())),
		format,
		args...,
	)
}

// Analyze performs sema analysis on the given AST program.
func (a *Analyzer) Analyze(program *ast.Program) error {
	a.pushScope()
	defer a.popScope()

	// Register every top level symbol first
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {

		case *ast.FuncStmt:
			a.declareFunction(&Symbol{
				Name:     n.Name.Value,
				Type:     n.ReturnType,
				Kind:     SymbolFunction,
				Function: n,
			})

		case *ast.StructStmt:
			a.declare(&Symbol{
				Name:   n.Name.Value,
				Kind:   SymbolStruct,
				Struct: n,
			})

		case *ast.EnumStmt:
			a.declare(&Symbol{
				Name: n.Name.Value,
				Kind: SymbolEnum,
				Enum: n,
			})

		}
	}

	a.assignFunctionCNames()

	// Infer function return types before sema analysis.
	//for _, stmt := range program.Statements {
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

			a.declare(&Symbol{
				Name:     param.Name.Value,
				Type:     param.Type,
				Kind:     SymbolVariable,
				DeclNode: param.Name,
			})
		}

		for _, stmt := range n.Body {
			a.visitStatement(stmt)
		}

		if n.ReturnType == nil {
			a.inferFunctionReturnType(n)

			if sym := a.lookupFunctionSymbol(n); sym != nil {
				sym.Type = n.ReturnType
			}
		}

		if n.ReturnType != nil && n.ReturnType.Value != "unit" {
			if !a.blockAllPathsReturn(n.Body) {
				a.errorf(
					n.Name,
					"function '%s' must return %s on all paths",
					n.Name.Value,
					n.ReturnType.Value,
				)
			}
		}

		a.popScope()

	case *ast.ReturnStmt:
		if a.currentFunction == nil {
			return
		}

		actual := &ast.Identifier{Value: "unit"}
		if n.Value != nil {
			actual = a.inferExprType(n.Value)
		}

		expected := a.currentFunction.ReturnType

		if expected != nil && actual != nil && expected.Value != actual.Value {
			var posErr ast.Node
			if n.Value == nil {
				posErr = n
			} else {
				posErr = n.Value
			}
			a.errorf(
				posErr,
				"return type mismatch: expected %s, got %s",
				expected.Value,
				actual.Value,
			)
		}

	case *ast.VarStmt:
		if n.Type == nil {
			n.Type = a.inferExprType(n.Value)
		} else if n.Value != nil {
			got := a.inferExprType(n.Value)

			if got != nil && got.Value != n.Type.Value {
				a.errorf(
					n.Value,
					"cannot initialize variable '%s' of type '%s' with value of type '%s'",
					n.Name.Value,
					n.Type.Value,
					got.Value,
				)
			}
		}

		if a.checkEnumShadow(n.Name) {
			return
		}

		a.declare(&Symbol{
			Name:     n.Name.Value,
			Type:     n.Type,
			Kind:     SymbolVariable,
			Mutable:  n.Mutable,
			DeclNode: n.Name,
		})

	case *ast.IfStmt:
		cond := a.inferExprType(n.Condition)
		if cond != nil && cond.Value != "bool" {
			a.errorf(n.Condition, "if condition must be bool, got %s", cond.Value)
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

			if got != nil && sym.Type != nil && got.Value != sym.Type.Value {
				a.errorf(
					n.Value,
					"cannot assign %s to variable '%s' of type %s",
					got.Value,
					left.Value,
					sym.Type.Value,
				)
			}

		case *ast.MemberExpr:
			objectType := a.inferExprType(left.Object)
			if objectType == nil {
				a.errorf(n.Value, "cannot determine type of member access")
				return
			}

			strct := a.lookupStruct(objectType.Value)
			if strct == nil {
				a.errorf(n.Value, "'%s' is not a struct", objectType.Value)
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

			if got != nil && got.Value != field.Type.Value {
				a.errorf(
					n.Value, "cannot assign %s to field '%s' of type %s",
					got.Value,
					field.Name.Value,
					field.Type.Value,
				)
			}

		default:
			a.errorf(n.Value, "left side of assignment is not assignable")
		}

	case *ast.ExpressionStmt:
		a.analyzeExpr(n.Expression)
	}
}

func (a *Analyzer) findReturnExprType(stmt ast.Stmt) *ast.Identifier {
	switch n := stmt.(type) {
	case *ast.ReturnStmt:
		if n.Value == nil {
			return &ast.Identifier{Value: "unit"}
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
	}

	return nil
}

func (a *Analyzer) blockAllPathsReturn(stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		if a.stmtAllPathsReturn(stmt) {
			return true
		}
	}
	return false
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

	default:
		return false
	}
}

func (a *Analyzer) inferFunctionReturnType(fn *ast.FuncStmt) {

	if fn.ReturnType != nil {
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
			fn.ReturnType = typ

			if sym != nil {
				sym.Type = typ
			}

			return
		}
	}

	fn.ReturnType = &ast.Identifier{Value: "unit"}

	if sym != nil {
		sym.Type = fn.ReturnType
	}
}

func (a *Analyzer) inferExprType(expr ast.Expr) *ast.Identifier {
	switch n := expr.(type) {

	case *ast.BadExpr:
		return nil

	case *ast.IntegerLiteral:
		return &ast.Identifier{Value: "int"}

	case *ast.FloatLiteral:
		return &ast.Identifier{Value: "float"}

	case *ast.CharacterLiteral:
		return &ast.Identifier{Value: "char"}

	case *ast.StringLiteral:
		return &ast.Identifier{Value: "string"}

	case *ast.BooleanLiteral:
		return &ast.Identifier{
			Value: "bool",
		}

	case *ast.Identifier:
		if sym := a.lookup(n.Value); sym != nil {
			return sym.Type
		}

		a.errorf(n, "unknown identifier: %s", n.Value)
		return nil

	case *ast.CallExpr:
		sym := a.resolveCallExpr(n)
		if sym == nil {
			return nil
		}

		if sym.Inferring {
			return nil
		}

		if sym.Type == nil {
			a.inferFunctionReturnType(sym.Function)
		}

		return sym.Type

	case *ast.BinaryExpr:
		left := a.inferExprType(n.Left)
		right := a.inferExprType(n.Right)

		if left == nil || right == nil {
			return nil
		}

		switch n.Operator.Kind {
		case token.Plus, token.Minus, token.Star, token.Slash:
			if (left.Value != "int" && left.Value != "float") ||
				(right.Value != "int" && right.Value != "float") {
				a.errorf(
					n,
					"operator '%s' requires numeric operands",
					n.TokenLiteral(),
				)
				return nil
			}

			if left.Value == "float" || right.Value == "float" {
				return &ast.Identifier{Value: "float"}
			}

			return &ast.Identifier{Value: "int"}

		case token.EqualEqual, token.BangEqual,
			token.Less, token.LessEqual,
			token.Greater, token.GreaterEqual:

			if left.Value != right.Value {
				a.errorf(
					n,
					"cannot compare %s with %s",
					left.Value,
					right.Value,
				)
				return nil
			}

			return &ast.Identifier{Value: "bool"}
		}

		return nil

	case *ast.MemberExpr:
		// the object is a type name used as a namespace, so it must be resolved before type inference
		if id, ok := n.Object.(*ast.Identifier); ok {
			if sym := a.lookup(id.Value); sym != nil && sym.Kind == SymbolEnum {
				for _, variant := range sym.Enum.Variants {
					if variant.Value == n.Property.Value {
						return &ast.Identifier{Value: sym.Enum.Name.Value}
					}
				}

				a.errorf(n.Property, "enum '%s' has no variant '%s'", sym.Enum.Name.Value, n.Property.Value)
				return nil
			}
		}

		objectType := a.inferExprType(n.Object)
		if objectType == nil {
			return nil
		}

		strct := a.lookupStruct(objectType.Value)
		if strct == nil {
			a.errorf(n.Object, "'%s' is not a struct", objectType.Value)
			return nil
		}

		for _, field := range strct.Fields {
			if field.Name.Value == n.Property.Value {
				return field.Type
			}
		}

		a.errorf(n.Property, "struct '%s' has no field '%s'", strct.Name.Value, n.Property.Value)
		return nil
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
