package sema

import (
	"github.com/azin-lang/Azin/pkg/ast"
	types2 "github.com/azin-lang/Azin/pkg/types"
)

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
				fieldType, err := a.lookupType(field.SynType.Value)
				if err != nil {
					a.errorf(field.SynType, "unknown type: %s", field.SynType.Value)
				}

				// Store the TypeInfo in the field node for later phase.
				field.SemaType = fieldType
			}

		case *ast.FuncStmt:
			returnType := types2.UnknownType()
			if n.SynReturnType != nil {
				foundType, err := a.lookupType(n.SynReturnType.Value)
				if err != nil {
					a.errorf(n.SynReturnType, "unknown return type: %s", n.SynReturnType.Value)
				}

				returnType = foundType
			}

			a.declareFunction(&Symbol{
				Name:     n.Name.Value,
				Type:     returnType,
				Kind:     SymbolFunction,
				Function: n,
			})

			for _, param := range n.Params {
				if param.SynType == nil {
					a.errorf(param.Name, "parameter '%s' must have a type", param.Name.Value)
					param.SemaType = types2.ErrorType()
					continue
				}

				paramType, err := a.lookupType(param.SynType.Value)
				if err != nil {
					a.errorf(param.SynType, "unknown parameter type: %s", param.SynType.Value)
				}

				// Store the TypeInfo in the parameter node for later phase.
				param.SemaType = paramType
			}

			n.SemaReturnType = returnType
		}
	}
}
