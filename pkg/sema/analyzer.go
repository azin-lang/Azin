package sema

import (
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/diagnostics"
)

/*
- The Analyzer struct is responsible for performing sema analysis on the AST of a program.
- It maintains a stack of scopes, keeps track of the current function being analyzed,
- and uses a diagnostics engine to report errors and warnings.
*/
type Analyzer struct {
	scopes []*Scope

	currentFunction *ast.FuncStmt
	diag            *diagnostics.Engine

	loopDepth int
}

// New creates a new Analyzer instance with the provided diagnostics engine.
func New(diag *diagnostics.Engine) *Analyzer {
	return &Analyzer{
		diag: diag,
	}
}
