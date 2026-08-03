package lexer_test

import (
	"testing"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/lexer"
	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

func TestLexer_Tokens(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []syntax.SyntaxKind
	}{
		{
			name:  "Basic Arithmetic",
			input: "10 + 20 * 5",
			expectedTokens: []syntax.SyntaxKind{
				syntax.IntegerLiteralToken,
				syntax.PlusToken,
				syntax.IntegerLiteralToken,
				syntax.StarToken,
				syntax.IntegerLiteralToken,
				syntax.EndOfFileToken,
			},
		},
		{
			name:  "Floats and scientific notation",
			input: "3.14 + 1e-5",
			expectedTokens: []syntax.SyntaxKind{
				syntax.FloatLiteralToken, // Ensure this exists in your syntax package
				syntax.PlusToken,
				syntax.FloatLiteralToken,
				syntax.EndOfFileToken,
			},
		},
		{
			name:  "Strings and Comments",
			input: `"hello" // say hi`,
			expectedTokens: []syntax.SyntaxKind{
				syntax.StringLiteralToken,
				syntax.EndOfFileToken,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := source.NewSourceText("test.az", 0, []byte(tt.input))
			diags := diagnostics.NewCollector()
			lex := lexer.New(file, diags)

			var tokens []syntax.SyntaxKind
			for {
				tok := lex.NextToken()
				tokens = append(tokens, tok.Kind())
				if tok.Kind() == syntax.EndOfFileToken {
					break
				}
			}

			if len(tokens) != len(tt.expectedTokens) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expectedTokens), len(tokens))
			}

			for i, expected := range tt.expectedTokens {
				if tokens[i] != expected {
					t.Errorf("token %d: expected %v, got %v", i, expected, tokens[i])
				}
			}
		})
	}
}
