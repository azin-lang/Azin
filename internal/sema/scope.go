package sema

// Scope represents a lexical scope in the source code. It contains a reference to its parent scope and a map of symbols defined within that scope.
type Scope struct {
	Parent    *Scope
	Symbols   map[string]*Symbol
	Functions map[string][]*Symbol
}

func (a *Analyzer) pushScope() {
	var parent *Scope
	if len(a.scopes) != 0 {
		parent = a.scopes[len(a.scopes)-1]
	}

	a.scopes = append(a.scopes, &Scope{
		Parent:    parent,
		Symbols:   map[string]*Symbol{},
		Functions: map[string][]*Symbol{},
	})
}

func (a *Analyzer) popScope() {
	scope := a.scopes[len(a.scopes)-1]
	for _, sym := range scope.Symbols {
		if sym.Kind == SymbolVariable && !sym.Used && sym.Name != "_" && sym.DeclNode != nil {
			a.warningf(sym.DeclNode, "unused variable: %s", sym.Name)
		}
	}
	a.scopes = a.scopes[:len(a.scopes)-1]
}

func (a *Analyzer) currentScope() *Scope {
	return a.scopes[len(a.scopes)-1]
}
