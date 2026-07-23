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
	parent    *Scope
	constants map[string]ast.Expr

	// modified tracks variables reassigned in this specific scope.
	// This is crucial for invalidating parent constants after branching.
	modified map[string]bool
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		currentScope: &Scope{
			constants: make(map[string]ast.Expr),
			modified:  make(map[string]bool),
		},
	}
}

// Optimize is the main entry point.
func Optimize(program *ast.Program) {
	if program == nil {
		return
	}
	opt := NewOptimizer()
	program.Statements = opt.optimizeStatements(program.Statements)
}

func (o *Optimizer) Enter() {
	o.currentScope = &Scope{
		parent:    o.currentScope,
		constants: make(map[string]ast.Expr),
		modified:  make(map[string]bool),
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

func (s *Scope) GetConstant(name string) (ast.Expr, bool) {
	if val, ok := s.constants[name]; ok {
		return val, true
	}
	if s.parent != nil {
		return s.parent.GetConstant(name)
	}
	return nil, false
}

func (s *Scope) SetConstant(name string, val ast.Expr) {
	s.constants[name] = val
	s.modified[name] = true
}

func (s *Scope) Invalidate(name string) {
	delete(s.constants, name)
	s.modified[name] = true
	if s.parent != nil {
		s.parent.Invalidate(name)
	}
}

func (s *Scope) ClearAll() {
	// Wipe this scope's constants
	s.constants = make(map[string]ast.Expr)

	// Recursively wipe parent scopes
	if s.parent != nil {
		s.parent.ClearAll()
	}
}
