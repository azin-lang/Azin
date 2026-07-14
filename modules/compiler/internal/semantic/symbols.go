package semantic

import "github.com/azin-lang/azin/compiler/internal/ast"

// SymbolKind represents the kind of a symbol (variable, function, struct).
type SymbolKind uint8

const (
	SymbolVariable SymbolKind = iota // SymbolVariable represents a variable symbol.
	SymbolFunction                   // SymbolFunction represents a function symbol.
	SymbolStruct                     // SymbolStruct represents a struct symbol.
)

// Symbol represents a symbol in the semantic analysis phase.
type Symbol struct {
	Name    string
	Type    *ast.Identifier
	Kind    SymbolKind
	Mutable bool

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
