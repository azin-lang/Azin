package semantic

type Scope struct {
	Parent  *Scope
	Symbols map[string]*Symbol
}

func (a *Analyzer) pushScope() {
	var parent *Scope
	if len(a.scopes) != 0 {
		parent = a.scopes[len(a.scopes)-1]
	}

	a.scopes = append(a.scopes, &Scope{
		Parent:  parent,
		Symbols: map[string]*Symbol{},
	})
}

func (a *Analyzer) popScope() {
	a.scopes = a.scopes[:len(a.scopes)-1]
}

func (a *Analyzer) currentScope() *Scope {
	return a.scopes[len(a.scopes)-1]
}
