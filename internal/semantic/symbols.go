package semantic

import "github.com/azin-lang/Azin/internal/ast"

type SymbolKind uint8

const (
	SymbolVariable SymbolKind = iota
	SymbolFunction
	SymbolStruct
)

type Symbol struct {
	Name string
	Type *ast.Identifier
	Kind SymbolKind

	Function *ast.FuncStmt
	Struct   *ast.StructStmt

	Inferring bool
}

func (a *Analyzer) lookup(name string) *Symbol {
	for scope := a.currentScope(); scope != nil; scope = scope.Parent {
		if sym, ok := scope.Symbols[name]; ok {
			return sym
		}
	}

	return nil
}

func (a *Analyzer) declare(sym *Symbol) {
	a.currentScope().Symbols[sym.Name] = sym
}
