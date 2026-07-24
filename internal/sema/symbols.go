package sema

import "github.com/azin-lang/Azin/internal/ast"

// SymbolKind represents the kind of a symbol (variable, function, struct).
type SymbolKind uint8

const (
	SymbolVariable SymbolKind = iota // SymbolVariable represents a variable symbol.
	SymbolFunction                   // SymbolFunction represents a function symbol.
	SymbolStruct                     // SymbolStruct represents a struct symbol.
	SymbolEnum                       // SymbolEnum represents an enum symbol.
)

// Symbol represents a symbol in the sema analysis phase.
type Symbol struct {
	Name     string
	Type     *ast.Identifier
	Kind     SymbolKind
	Mutable  bool
	Used     bool
	DeclNode ast.Node

	Function *ast.FuncStmt
	Struct   *ast.StructStmt
	Enum     *ast.EnumStmt

	Inferring bool
}

func (a *Analyzer) lookup(name string) *Symbol {
	for scope := a.currentScope(); scope != nil; scope = scope.Parent {
		if sym, ok := scope.Symbols[name]; ok {
			if sym.Kind == SymbolVariable {
				sym.Used = true
			}
			return sym
		}
	}

	return nil
}

func (a *Analyzer) lookupFunctions(name string) []*Symbol {
	for scope := a.currentScope(); scope != nil; scope = scope.Parent {
		if syms, ok := scope.Functions[name]; ok {
			return syms
		}
	}

	return nil
}

func (a *Analyzer) declare(sym *Symbol) {
	if existing := a.currentScope().Symbols[sym.Name]; existing != nil {
		a.errorfSym(sym, "redeclaration of '%s'", sym.Name)
		return
	}
	if existing := a.currentScope().Functions[sym.Name]; len(existing) != 0 {
		a.errorfSym(sym, "redeclaration of '%s'", sym.Name)
		return
	}
	a.currentScope().Symbols[sym.Name] = sym
}

func sameParamTypes(left, right *ast.FuncStmt) bool {
	if len(left.Params) != len(right.Params) {
		return false
	}

	for i := range left.Params {
		if left.Params[i].Type == nil || right.Params[i].Type == nil {
			return false
		}

		if left.Params[i].Type.Value != right.Params[i].Type.Value {
			return false
		}
	}

	return true
}

func (a *Analyzer) declareFunction(sym *Symbol) {
	if existing := a.currentScope().Symbols[sym.Name]; existing != nil {
		a.errorfSym(sym, "redeclaration of '%s'", sym.Name)
		return
	}

	for _, existing := range a.currentScope().Functions[sym.Name] {
		if sameParamTypes(existing.Function, sym.Function) {
			a.errorfSym(sym, "redeclaration of '%s'", sym.Name)
			return
		}
	}

	a.currentScope().Functions[sym.Name] = append(a.currentScope().Functions[sym.Name], sym)
}

func (a *Analyzer) lookupFunctionSymbol(fn *ast.FuncStmt) *Symbol {
	for scope := a.currentScope(); scope != nil; scope = scope.Parent {
		for _, sym := range scope.Functions[fn.Name.Value] {
			if sym.Function == fn {
				return sym
			}
		}
	}

	return nil
}

//nolint:unparam
func (a *Analyzer) errorfSym(sym *Symbol, format string, args ...any) {
	switch {
	case sym.Function != nil:
		a.errorf(sym.Function, format, args...)
	case sym.Struct != nil:
		a.errorf(sym.Struct, format, args...)
	case sym.Enum != nil:
		a.errorf(sym.Enum, format, args...)
	}
}
