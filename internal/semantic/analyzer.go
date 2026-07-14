package semantic

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/diagnostics"
)

/*
- The Analyzer struct is responsible for performing semantic analysis on the AST of a program.
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

func (a *Analyzer) errorf(node ast.Node, format string, args ...any) {
	a.diag.ReportError(
		node.Pos(),
		uint32(len(node.TokenLiteral())),
		format,
		args...,
	)
}

// Analyze performs semantic analysis on the given AST program.
func (a *Analyzer) Analyze(program *ast.Program) error {
	a.pushScope()
	defer a.popScope()

	// Register every top level symbol first
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {

		case *ast.FuncStmt:
			a.declare(&Symbol{
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

	// Analyze every statement.
	for _, stmt := range program.Statements {
		a.visitStatement(stmt)
	}

	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncStmt); ok {
			a.inferFunctionReturnType(fn)

			sym := a.lookup(fn.Name.Value)
			sym.Type = fn.ReturnType
		}
	}

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
				Name: param.Name.Value,
				Type: param.Type,
				Kind: SymbolVariable,
			})
		}

		for _, stmt := range n.Body {
			a.visitStatement(stmt)
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
			a.errorf(
				n.Value,
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
			Name:    n.Name.Value,
			Type:    n.Type,
			Kind:    SymbolVariable,
			Mutable: n.Mutable,
		})

	case *ast.IfStmt:
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
	}
}

func (a *Analyzer) findReturnType(stmt ast.Stmt) *ast.Identifier {
	switch n := stmt.(type) {
	case *ast.BadStmt:
		return nil

	case *ast.ReturnStmt:
		if n.Value == nil {
			return &ast.Identifier{Value: "unit"}
		}
		return a.inferExprType(n.Value)

	case *ast.IfStmt:
		for _, s := range n.Then {
			if t := a.findReturnType(s); t != nil {
				return t
			}
		}

		for _, s := range n.Else {
			if t := a.findReturnType(s); t != nil {
				return t
			}
		}
	}

	return nil
}

func (a *Analyzer) inferFunctionReturnType(fn *ast.FuncStmt) {

	if fn.ReturnType != nil {
		return
	}

	sym := a.lookup(fn.Name.Value)
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
		if typ := a.findReturnType(stmt); typ != nil {
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

	case *ast.CallExpr:
		id, ok := n.Callee.(*ast.Identifier)
		if !ok {
			return nil
		}

		sym := a.lookup(id.Value)
		if sym == nil || sym.Kind != SymbolFunction {
			a.errorf(n.Callee, "unknown function: %s", id.Value)
			return nil
		}

		if sym.Inferring {
			return nil
		}

		if sym.Type == nil {
			a.inferFunctionReturnType(sym.Function)
		}

		if len(n.Args) != len(sym.Function.Params) {
			a.errorf(n.Callee, "wrong number of arguments to %s", id.Value)
			return nil
		}

		for i, arg := range n.Args {
			got := a.inferExprType(arg)
			want := sym.Function.Params[i].Type

			if got == nil {
				continue
			}

			if got.Value != want.Value {
				a.errorf(
					arg,
					"argument %d of %s: expected %s, got %s",
					i+1,
					id.Value,
					want.Value,
					got.Value,
				)
			}
		}

		return sym.Type

	case *ast.BinaryExpr:
		left := a.inferExprType(n.Left)
		right := a.inferExprType(n.Right)

		if left == nil || right == nil {
			return nil
		}

		if left.Value == "float" || right.Value == "float" {
			return &ast.Identifier{Value: "float"}
		}

		return left

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
