package compiler

import (
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/source"
)

// compiles the given file.
func Compile(file *source.File) error {
	diag := diagnostics.New(file)
	_ = lexer.New(file, diag).Tokenize()

	return diag.Err()
}
