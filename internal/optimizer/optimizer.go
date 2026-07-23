// Package optimizer provides basic optimizations before transpiling to C.
package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

// Optimizer holds the global state for the optimization pass.
type Optimizer struct {
	currentScope *Scope
}

// Scope tracks variables for a specific block (e.g., global, function, if-body).
type Scope struct {
	parent *Scope
	values map[string]ast.Expr

	// modified tracks variables reassigned in this specific scope.
	// This is crucial for invalidating parent values after branching.
	modified map[string]bool
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		currentScope: &Scope{
			values:   make(map[string]ast.Expr),
			modified: make(map[string]bool),
		},
	}
}

// Optimize is the main entry point.
func Optimize(program *ast.Program) {
	if program == nil {
		return
	}
	opt := NewOptimizer()

	// Pre-scan statements to register enum constants into the global scope
	for _, stmt := range program.Statements {
		if enumStmt, ok := stmt.(*ast.EnumStmt); ok {
			for i, field := range enumStmt.Variants {
				// Register e.g., "Color.Red" = 0, "Color.Green" = 1, etc.
				key := enumStmt.Name.Value + "." + field.Value
				opt.currentScope.SetValue(key, intLit(int64(i)))
			}
		}
	}

	program.Statements = opt.optimizeStatements(program.Statements)
}

func (o *Optimizer) Enter() {
	o.currentScope = &Scope{
		parent:   o.currentScope,
		values:   make(map[string]ast.Expr),
		modified: make(map[string]bool),
	}
}

func (o *Optimizer) Leave() {
	child := o.currentScope
	o.currentScope = child.parent

	if o.currentScope != nil {
		for name := range child.modified {
			o.currentScope.Invalidate(name)
		}
	}
}

func (s *Scope) GetValue(name string) (ast.Expr, bool) {
	if val, ok := s.values[name]; ok {
		return val, true
	}
	if s.parent != nil {
		return s.parent.GetValue(name)
	}
	return nil, false
}

func (s *Scope) SetValue(name string, val ast.Expr) {
	s.values[name] = val
	s.modified[name] = true
}

func (s *Scope) Invalidate(name string) {
	delete(s.values, name)
	s.modified[name] = true

	// Invalidate any variable holding a propagated copy of this variable
	// We extract keys to safely modify the map during iteration
	var aliases []string
	for k, v := range s.values {
		if id, ok := v.(*ast.Identifier); ok && id.Value == name {
			aliases = append(aliases, k)
		}
	}

	for _, alias := range aliases {
		s.Invalidate(alias)
	}

	if s.parent != nil {
		s.parent.Invalidate(name)
	}
}

func (s *Scope) ClearAll() {
	// Wipe this scope's known values
	s.values = make(map[string]ast.Expr)

	// Recursively wipe parent scopes
	if s.parent != nil {
		s.parent.ClearAll()
	}
}
