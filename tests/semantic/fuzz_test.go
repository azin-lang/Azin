package semantic_test

import (
	"testing"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/sema"
	"github.com/azin-lang/Azin/internal/source"
)

func FuzzSemantic(f *testing.F) {
	seeds := []string{
		"fn main: int do return 0; end",
		"fn foo: int do var x: int = 42; return x; end",
		"fn main: int do var x: int = 42; x = 99; return x; end",
		"fn add(a: int, b: int): int do return a + b; end",
		"struct Point is x: int; end fn main: int do var p: Point; return p.x; end",
		"fn foo do return 42; end",
		"importc \"stdio.h\"",
		"",
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		file := source.New("fuzz.az", data)
		diag := diagnostics.New(file)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic: %v", r)
			}
		}()

		tokens := lexer.New(file, diag).Tokenize()
		program, err := parser.Parse(string(file.Slice(0, file.Len())), tokens, diag)
		if err != nil {
			return
		}

		analyzer := sema.New(diag)
		analyzer.Analyze(program)
	})
}
