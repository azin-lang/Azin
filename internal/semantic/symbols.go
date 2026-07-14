package semantic

import "github.com/azin-lang/Azin/internal/ast"

// SymbolKind represents the kind of a symbol (variable, function, struct).
type SymbolKind uint8

const (
	SymbolVariable SymbolKind = iota // SymbolVariable represents a variable symbol.
	SymbolFunction                   // SymbolFunction represents a function symbol.
	SymbolStruct                     // SymbolStruct represents a struct symbol.
	SymbolEnum                       // SymbolEnum represents an enum symbol.
)

// Symbol represents a symbol in the semantic analysis phase.
type Symbol struct {
	Name    string
	Type    *ast.Identifier
	Kind    SymbolKind
	Mutable bool

	Function *ast.FuncStmt
	Struct   *ast.StructStmt
	Enum     *ast.EnumStmt

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
	if existing := a.currentScope().Symbols[sym.Name]; existing != nil {
		a.errorfSym(sym, "redeclaration of '%s'", sym.Name)
		return
	}
	a.currentScope().Symbols[sym.Name] = sym
}

func (a *Analyzer) errorfSym(sym *Symbol, format string, args ...any) {
	if sym.Function != nil {
		a.errorf(sym.Function, format, args...)
	} else if sym.Struct != nil {
		a.errorf(sym.Struct, format, args...)
	} else if sym.Enum != nil {
		a.errorf(sym.Enum, format, args...)
	}
}
