// Package optimizer provides basic optimizations before transpiling to C.
package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
)

// Optimize traverses the program AST and applies compile-time simplifications.
// It currently performs constant folding and conditional dead-code elimination.
func Optimize(program *ast.Program) {
	if program == nil {
		return
	}

	program.Statements = optimizeStatements(program.Statements)
}
