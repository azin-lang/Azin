package sema

import (
	"fmt"

	"github.com/azin-lang/Azin/pkg/ast"
	types2 "github.com/azin-lang/Azin/pkg/types"
)

func (a *Analyzer) lookupType(name string) (*types2.TypeInfo, error) {
	if primitive := types2.GetPrimitiveTypeByName(name); primitive != nil {
		return primitive, nil
	}

	sym := a.lookup(name)
	if sym == nil {
		return types2.ErrorType(), fmt.Errorf("unknown type: %s", name)
	}

	switch sym.Kind {
	case SymbolStruct, SymbolEnum:
		return types2.NominalType(sym.Name), nil
	default:
		if sym.Type != nil {
			return sym.Type, nil
		}

		return types2.ErrorType(), fmt.Errorf("internal compiler error: symbol has no type info: %s", sym.Name)
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
