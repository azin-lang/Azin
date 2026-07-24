package parser_test

import (
	"testing"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/source"
)

func FuzzParser(f *testing.F) {
	seeds := []string{
		"var x: int = 42;",
		"fn main: int do return 0; end",
		"if true then return 1; end",
		"if true then return 1; else return 2; end",
		"loop return 0; end",
		"struct Point is x: int; y: int; end",
		"importc \"stdio.h\"",
		"x = 42;",
		"foo();",
		"foo(1, 2, 3);",
		"point.x;",
		"fn foo do return; end",
		"var mut x: int = 42;",
		"var s: string = \"hello\";",
		"var c: char = 'x';",
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
		_, _ = parser.Parse(string(file.Slice(0, file.Len())), tokens, diag)
	})
}
