package semantic

import (
	"strconv"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/diagnostics"
)

type Analyzer struct {
	scopes []*Scope
	diag   *diagnostics.Engine

	currentFunction *ast.FuncStmt
}

func New(diag *diagnostics.Engine) *Analyzer {
	return &Analyzer{diag: diag}
}

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

func (a *Analyzer) visitStatement(stmt ast.Stmt) {
	switch n := stmt.(type) {

	case *ast.FuncStmt:
		old := a.currentFunction
		a.currentFunction = n
		defer func() {
			a.currentFunction = old
		}()

		a.pushScope()

		// Register parameters.
		for _, param := range n.Params {
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
			panic(
				"return type mismatch: expected " +
					expected.Value +
					", got " +
					actual.Value,
			)
		}

	case *ast.VarStmt:
		if n.Type == nil {
			n.Type = a.inferExprType(n.Value)
		} else if n.Value != nil {
			got := a.inferExprType(n.Value)

			if got != nil && got.Value != n.Type.Value {
				panic(
					"cannot initialize variable '" +
						n.Name.Value +
						"' of type " +
						n.Type.Value +
						"' with value of type " +
						got.Value,
				)
			}
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

	case *ast.AssignmentStmt:
		switch left := n.Left.(type) {

		case *ast.Identifier:
			sym := a.lookup(left.Value)
			if sym == nil {
				panic("unknown variable: " + left.Value)
			}

			if sym.Kind != SymbolVariable {
				panic(left.Value + " is not a variable")
			}

			if !sym.Mutable {
				panic("cannot assign to immutable variable '" + left.Value + "'")
			}

			got := a.inferExprType(n.Value)

			if got != nil && sym.Type != nil && got.Value != sym.Type.Value {
				panic(
					"cannot assign " +
						got.Value +
						" to variable '" +
						left.Value +
						"' of type " +
						sym.Type.Value,
				)
			}

		case *ast.MemberExpr:
			objectType := a.inferExprType(left.Object)
			if objectType == nil {
				panic("cannot determine type of member access")
			}

			strct := a.lookupStruct(objectType.Value)
			if strct == nil {
				panic("'" + objectType.Value + "' is not a struct")
			}

			field := a.lookupField(strct, left.Property.Value)
			if field == nil {
				panic(
					"struct '" +
						strct.Name.Value +
						"' has no field '" +
						left.Property.Value +
						"'",
				)
			}

			if !field.Mutable {
				panic(
					"cannot assign to immutable field '" +
						field.Name.Value +
						"'",
				)
			}

			got := a.inferExprType(n.Value)

			if got != nil && got.Value != field.Type.Value {
				panic(
					"cannot assign " +
						got.Value +
						" to field '" +
						field.Name.Value +
						"' of type " +
						field.Type.Value,
				)
			}

		default:
			panic("left side of assignment is not assignable")
		}
	}
}

func (a *Analyzer) findReturnType(stmt ast.Stmt) *ast.Identifier {
	switch n := stmt.(type) {

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

	case *ast.IntegerLiteral:
		return &ast.Identifier{Value: "int"}

	case *ast.FloatLiteral:
		return &ast.Identifier{Value: "float"}

	case *ast.CharacterLiteral:
		return &ast.Identifier{Value: "char"}

	case *ast.StringLiteral:
		return &ast.Identifier{Value: "string"}

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
			panic("unknown function: " + id.Value)
		}

		if sym.Inferring {
			return nil
		}

		if sym.Type == nil {
			a.inferFunctionReturnType(sym.Function)
		}

		if len(n.Args) != len(sym.Function.Params) {
			panic("wrong number of arguments to " + id.Value)
		}

		for i, arg := range n.Args {
			got := a.inferExprType(arg)
			want := sym.Function.Params[i].Type

			if got == nil {
				continue
			}

			if got.Value != want.Value {
				panic(
					"argument " +
						strconv.Itoa(i+1) +
						" of " +
						id.Value +
						": expected " +
						want.Value +
						", got " +
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
		// Enum.Variant, the object is a type name used as a namespace, so it must be resolved before type inference
		if id, ok := n.Object.(*ast.Identifier); ok {
			if sym := a.lookup(id.Value); sym != nil && sym.Kind == SymbolEnum {
				for _, v := range sym.Enum.Variants {
					if v.Value == n.Property.Value {
						return &ast.Identifier{Value: sym.Enum.Name.Value}
					}
				}
				a.diag.ReportError(
					n.Property.Token.Position,
					n.Property.Token.Length,
					"enum '%s' has no variant '%s'",
					sym.Enum.Name.Value,
					n.Property.Value,
				)
				return nil
			}
		}

		objectType := a.inferExprType(n.Object)
		if objectType == nil {
			return nil
		}

		strct := a.lookupStruct(objectType.Value)
		if strct == nil {
			panic("'" + objectType.Value + "' is not a struct")
		}

		for _, field := range strct.Fields {
			if field.Name.Value == n.Property.Value {
				return field.Type
			}
		}

		panic(
			"struct '" +
				strct.Name.Value +
				"' has no field '" +
				n.Property.Value +
				"'",
		)
	}

	return nil
}
