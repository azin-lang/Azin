package lexer_test

import (
	"testing"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/source"
)

func FuzzLexer(f *testing.F) {
	seeds := []string{
		"var x: int = 42;",
		"fn main: int do return 0; end",
		"\"hello\"",
		"'x'",
		"// comment\n42",
		"/* block */ 42",
		"42",
		"3.14",
		"+ - * / % = == ! !=",
		"< <= > >= += ++ -= -- -> && ||",
		"( ) { } [ ] , ; : .",
		"fn do var mut return end char int bool unit string float if then else struct is importc loop null",
		"@",
		"'\n'",
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

		lexer.New(file, diag).Tokenize()
	})
}
