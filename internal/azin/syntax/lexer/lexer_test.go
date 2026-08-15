package lexer_test

import (
	"testing"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/syntax/lexer"
	"github.com/azin-lang/Azin/internal/azin/text"
)

// Helper function to lex an entire input string into green tokens
func lexAll(t *testing.T, input string) ([]*green.Token, *diagnostics.Collector) {
	t.Helper()

	file := text.NewSourceText("test.az", 0, []byte(input))
	diags := diagnostics.NewCollector()
	lex := lexer.New(file, diags)

	var tokens []*green.Token
	for {
		tok := lex.NextToken()
		tokens = append(tokens, tok)
		if tok.Kind() == syntax.EndOfFileToken {
			break
		}
	}

	return tokens, diags
}

func TestLexer_Tokens(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []syntax.Kind
	}{
		{
			name:  "Basic Arithmetic",
			input: "10 + 20 * 5",
			expectedTokens: []syntax.Kind{
				syntax.IntegerLiteralToken,
				syntax.PlusToken,
				syntax.IntegerLiteralToken,
				syntax.StarToken,
				syntax.IntegerLiteralToken,
				syntax.EndOfFileToken,
			},
		},
		{
			name:  "Floats and Scientific Notation",
			input: "3.14 + 1e-5",
			expectedTokens: []syntax.Kind{
				syntax.FloatLiteralToken,
				syntax.PlusToken,
				syntax.FloatLiteralToken,
				syntax.EndOfFileToken,
			},
		},
		{
			name:  "Strings and Comments",
			input: `"hello" // say hi`,
			expectedTokens: []syntax.Kind{
				syntax.StringLiteralToken,
				syntax.EndOfFileToken,
			},
		},
		{
			name:  "Delimiters and Grouping",
			input: "(a, b) { c; }",
			expectedTokens: []syntax.Kind{
				syntax.OpenParenToken,
				syntax.IdentifierToken,
				syntax.CommaToken,
				syntax.IdentifierToken,
				syntax.CloseParenToken,
				syntax.OpenBraceToken,
				syntax.IdentifierToken,
				syntax.SemicolonToken,
				syntax.CloseBraceToken,
				syntax.EndOfFileToken,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := lexAll(t, tt.input)

			if len(tokens) != len(tt.expectedTokens) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expectedTokens), len(tokens))
			}

			for i, expected := range tt.expectedTokens {
				if tokens[i].Kind() != expected {
					t.Errorf("token %d: expected %v, got %v", i, expected, tokens[i].Kind())
				}
			}
		})
	}
}

func TestLexer_TokenText(t *testing.T) {
	tests := []struct {
		input        string
		expectedText []string
	}{
		{
			input:        "var score = 100",
			expectedText: []string{"var", "score", "=", "100", ""},
		},
		{
			input:        `"hello world"`,
			expectedText: []string{`"hello world"`, ""},
		},
	}

	for _, tt := range tests {
		tokens, _ := lexAll(t, tt.input)

		if len(tokens) != len(tt.expectedText) {
			t.Fatalf("expected %d tokens, got %d", len(tt.expectedText), len(tokens))
		}

		for i, expected := range tt.expectedText {
			if tokens[i].Text() != expected {
				t.Errorf("token %d: expected text %q, got %q", i, expected, tokens[i].Text())
			}
		}
	}
}

func TestLexer_MultiCharacterOperators(t *testing.T) {
	input := "== != <= >= && || += -="
	expected := []syntax.Kind{
		syntax.EqualsEqualsToken,
		syntax.BangEqualsToken,
		syntax.LessEqualsToken,
		syntax.GreaterEqualsToken,
		syntax.AmpersandAmpersandToken,
		syntax.PipePipeToken,
		syntax.PlusEqualsToken,
		syntax.MinusEqualsToken,
		syntax.EndOfFileToken,
	}

	tokens, diags := lexAll(t, input)
	if diags.HasErrors() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Kind() != exp {
			t.Errorf("token %d: expected %v, got %v", i, exp, tokens[i].Kind())
		}
	}
}

func TestLexer_KeywordsVsIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected syntax.Kind
	}{
		{"var", syntax.KeywordVar},
		{"fn", syntax.KeywordFn},
		{"if", syntax.KeywordIf},
		{"else", syntax.KeywordElse},
		{"return", syntax.KeywordReturn},
		{"true", syntax.KeywordTrue},
		{"false", syntax.KeywordFalse},
		{"letter", syntax.IdentifierToken},
		{"fnName", syntax.IdentifierToken},
		{"_private", syntax.IdentifierToken},
		{"var123", syntax.IdentifierToken},
	}

	for _, tt := range tests {
		tokens, _ := lexAll(t, tt.input)
		if tokens[0].Kind() != tt.expected {
			t.Errorf("for input %q: expected %v, got %v", tt.input, tt.expected, tokens[0].Kind())
		}
	}
}

func TestLexer_TriviaPreservation(t *testing.T) {
	// Leading spaces attached to '10', trailing comment attached to EOF or '20'
	input := "  10 + /* inline */ 20 // line comment"
	tokens, _ := lexAll(t, input)

	// First token '10' should have leading trivia (2 spaces)
	tok10 := tokens[0]
	if tok10.Text() != "10" {
		t.Fatalf("expected '10', got %q", tok10.Text())
	}

	if tok10.LeadingTrivia() == nil {
		t.Fatal("expected leading trivia for token '10'")
	}

	// Total full width should match raw input string length
	var totalWidth int
	for _, tok := range tokens {
		totalWidth += int(tok.FullWidth())
	}

	if totalWidth != len(input) {
		t.Fatalf("expected full width %d, got %d", len(input), totalWidth)
	}
}

func TestLexer_UnclosedStringDiagnostic(t *testing.T) {
	input := `"unclosed string literal`
	tokens, diags := lexAll(t, input)

	if !diags.HasErrors() {
		t.Fatal("expected diagnostic error for unclosed string")
	}

	if len(tokens) == 0 {
		t.Fatal("lexer should return at least BadToken or StringLiteralToken")
	}
}

func TestLexer_InvalidCharacterDiagnostic(t *testing.T) {
	input := `@`
	_, diags := lexAll(t, input)

	if !diags.HasErrors() {
		t.Fatal("expected diagnostic error for invalid character '@'")
	}
}
