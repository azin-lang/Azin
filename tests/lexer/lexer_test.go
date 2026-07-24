//nolint:goconst
package lexer_test

// If the lexer crashes, it's not our fault.
// You probably used a character that doesn't exist yet.
import (
	"strings"
	"testing"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

func lex(input string) ([]token.Token, *diagnostics.Engine) {
	file := source.New("test.az", []byte(input))
	diag := diagnostics.New(file)
	tokens := lexer.New(file, diag).Tokenize()
	return tokens, diag
}

func kindString(tok token.Token) string {
	return tok.Kind.String()
}

func joinKinds(tokens []token.Token) string {
	kinds := make([]string, 0, len(tokens))
	for _, t := range tokens {
		kinds = append(kinds, t.Kind.String())
	}
	return strings.Join(kinds, " ")
}

func TestLexerKeywords(t *testing.T) {
	input := "fn do var mut return end char int bool unit string float if then else struct is importc loop null"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	got := joinKinds(tokens)
	want := "kw_fn kw_do kw_var kw_mut kw_return kw_end kw_char kw_int kw_bool kw_unit kw_string kw_float kw_if kw_then kw_else kw_struct kw_is kw_import kw_loop kw_null eof"

	if got != want {
		t.Errorf("keywords\ngot:  %s\nwant: %s", got, want)
	}
}

func TestLexerIdentifiers(t *testing.T) {
	input := "foo bar _baz _ hello_world"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	idents := 0
	for _, tok := range tokens {
		if tok.Kind == token.Identifier {
			idents++
		}
	}
	if idents != 5 {
		t.Errorf("expected 5 identifiers, got %d", idents)
	}
}

func TestLexerIntegers(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"42", "decimal"},
		{"0", "zero"},
		{"0xFF", "hex"},
		{"0xDEAD_BEEF", "hex_underscore"},
		{"0b1010", "binary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, diag := lex(tt.input)
			if diag.HasErrors() {
				t.Fatalf("unexpected errors: %v", diag.Err())
			}
			if len(tokens) < 1 || tokens[0].Kind != token.IntegerLiteral {
				t.Errorf("expected IntegerLiteral for %s, got %s", tt.input, tokens[0].Kind)
			}
		})
	}
}

func TestLexerFloats(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"3.14", "simple"},
		{"0.5", "leading_zero"},
		{"10.0", "trailing_zero"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, diag := lex(tt.input)
			if diag.HasErrors() {
				t.Fatalf("unexpected errors: %v", diag.Err())
			}
			if len(tokens) < 1 || tokens[0].Kind != token.FloatLiteral {
				t.Errorf("expected FloatLiteral for %s, got %s", tt.input, tokens[0].Kind)
			}
		})
	}
}

func TestLexerStrings(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{`"hello"`, "simple"},
		{`""`, "empty"},
		{`"hello\nworld"`, "escape_newline"},
		{`"tab\there"`, "escape_tab"},
		{`"quote\"inside"`, "escape_quote"},
		{`"backslash\\here"`, "escape_backslash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, diag := lex(tt.input)
			if diag.HasErrors() {
				t.Fatalf("unexpected errors: %v", diag.Err())
			}
			if len(tokens) < 1 || tokens[0].Kind != token.StringLiteral {
				t.Errorf("expected StringLiteral for %s, got %s", tt.input, tokens[0].Kind)
			}
		})
	}
}

func TestLexerChars(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"'x'", "simple"},
		{`'\n'`, "escape_newline"},
		{`'\t'`, "escape_tab"},
		{`'\\'`, "escape_backslash"},
		{`'\''`, "escape_quote"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, diag := lex(tt.input)
			if diag.HasErrors() {
				t.Fatalf("unexpected errors: %v", diag.Err())
			}
			if len(tokens) < 1 || tokens[0].Kind != token.CharacterLiteral {
				t.Errorf("expected CharacterLiteral for %s, got %s", tt.input, tokens[0].Kind)
			}
		})
	}
}

func TestLexerOperators(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plus", "+", "plus"},
		{"minus", "-", "minus"},
		{"star", "*", "star"},
		{"slash", "/", "slash"},
		{"modulo", "%", "modulo"},
		{"equal", "=", "equal"},
		{"equal_equal", "==", "equal_equal"},
		{"bang", "!", "bang"},
		{"bang_equal", "!=", "bang_equal"},
		{"less", "<", "less"},
		{"less_equal", "<=", "less_equal"},
		{"greater", ">", "greater"},
		{"greater_equal", ">=", "greater_equal"},
		{"plus_equal", "+=", "plus_equal"},
		{"plus_plus", "++", "plus_plus"},
		{"minus_equal", "-=", "minus_equal"},
		{"minus_minus", "--", "minus_minus"},
		{"arrow", "->", "arrow"},
		{"logical_and", "&&", "logical_and"},
		{"logical_or", "||", "logical_or"},
		{"ampersand", "&", "ampersand"},
		{"pipe", "|", "pipe"},
		{"less_less", "<<", "less_less"},
		{"greater_greater", ">>", "greater_greater"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, diag := lex(tt.input)
			if diag.HasErrors() {
				t.Fatalf("unexpected errors: %v", diag.Err())
			}
			if len(tokens) < 1 || kindString(tokens[0]) != tt.want {
				t.Errorf("expected %s for %q, got %s", tt.want, tt.input, kindString(tokens[0]))
			}
		})
	}
}

func TestLexerPunctuation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"left_paren", "(", "left_paren"},
		{"right_paren", ")", "right_paren"},
		{"left_brace", "{", "left_brace"},
		{"right_brace", "}", "right_brace"},
		{"left_bracket", "[", "left_bracket"},
		{"right_bracket", "]", "right_bracket"},
		{"comma", ",", "comma"},
		{"semicolon", ";", "semicolon"},
		{"colon", ":", "colon"},
		{"dot", ".", "dot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, diag := lex(tt.input)
			if diag.HasErrors() {
				t.Fatalf("unexpected errors: %v", diag.Err())
			}
			if len(tokens) < 1 || kindString(tokens[0]) != tt.want {
				t.Errorf("expected %s for %q, got %s", tt.want, tt.input, kindString(tokens[0]))
			}
		})
	}
}

func TestLexerLineComment(t *testing.T) {
	input := "// this is a comment\n42"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	// Should have newline, integer_literal, eof
	hasInt := false
	for _, tok := range tokens {
		if tok.Kind == token.IntegerLiteral {
			hasInt = true
			break
		}
	}
	if !hasInt {
		t.Errorf("expected IntegerLiteral after comment, got tokens: %s", joinKinds(tokens))
	}
}

func TestLexerBlockComment(t *testing.T) {
	input := "/* block comment */ 42"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	hasInt := false
	for _, tok := range tokens {
		if tok.Kind == token.IntegerLiteral {
			hasInt = true
			break
		}
	}
	if !hasInt {
		t.Errorf("expected IntegerLiteral after block comment, got tokens: %s", joinKinds(tokens))
	}
}

func TestLexerNestedBlockComment(t *testing.T) {
	input := "/* outer /* inner */ */ 99"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	hasInt := false
	for _, tok := range tokens {
		if tok.Kind == token.IntegerLiteral {
			hasInt = true
			break
		}
	}
	if !hasInt {
		t.Errorf("expected IntegerLiteral after nested block comment, got tokens: %s", joinKinds(tokens))
	}
}

func TestLexerUnterminatedString(t *testing.T) {
	input := `"hello`
	tokens, diag := lex(input)

	if !diag.HasErrors() {
		t.Error("expected error for unterminated string")
	}
	if len(tokens) < 1 || tokens[0].Kind != token.StringLiteral {
		t.Errorf("expected StringLiteral token, got %s", kindString(tokens[0]))
	}
}

func TestLexerEmptyChar(t *testing.T) {
	input := "''"
	tokens, diag := lex(input)

	if !diag.HasErrors() {
		t.Error("expected error for empty char literal")
	}
	if len(tokens) < 1 || tokens[0].Kind != token.CharacterLiteral {
		t.Errorf("expected CharacterLiteral token, got %s", kindString(tokens[0]))
	}
}

func TestLexerNewlines(t *testing.T) {
	input := "a\nb\n\nc"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	newlineCount := 0
	for _, tok := range tokens {
		if tok.Kind == token.Newline {
			newlineCount++
		}
	}
	// a\n = newline, b\n = newline, \n = newline, c = no newline
	if newlineCount != 3 {
		t.Errorf("expected 3 newlines, got %d. Tokens: %s", newlineCount, joinKinds(tokens))
	}
}

func TestLexerCRLF(t *testing.T) {
	input := "a\r\nb"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	if len(tokens) < 3 || tokens[0].Kind != token.Identifier ||
		tokens[1].Kind != token.Newline ||
		tokens[2].Kind != token.Identifier {
		t.Errorf("unexpected tokens for CRLF: %s", joinKinds(tokens))
	}
}

func TestLexerUnterminatedBlockComment(t *testing.T) {
	input := "/* unclosed"
	tokens, diag := lex(input)

	if !diag.HasErrors() {
		t.Error("expected error for unterminated block comment")
	}
	_ = tokens
}

func TestLexerInvalidEscapeString(t *testing.T) {
	input := `"\z"`
	tokens, diag := lex(input)

	if !diag.HasErrors() {
		t.Error("expected error for invalid escape in string")
	}
	if len(tokens) < 1 || tokens[0].Kind != token.StringLiteral {
		t.Errorf("expected StringLiteral token, got %s", kindString(tokens[0]))
	}
}

func TestLexerUnknownChar(t *testing.T) {
	input := "@"
	tokens, diag := lex(input)

	if !diag.HasErrors() {
		t.Error("expected error for unknown character")
	}
	if len(tokens) < 2 || tokens[0].Kind != token.Unknown {
		t.Errorf("expected Unknown token, got %s", kindString(tokens[0]))
	}
}

func TestLexerMultipleTokens(t *testing.T) {
	input := "fn add(a: int, b: int): int do\n    return a + b;\nend"
	tokens, diag := lex(input)

	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}

	expectedKinds := []token.Kind{
		token.KwFn, token.Identifier,
		token.LeftParen, token.Identifier, token.Colon, token.KwInt,
		token.Comma, token.Identifier, token.Colon, token.KwInt,
		token.RightParen, token.Colon, token.KwInt, token.KwDo,
		token.Newline,
		token.KwReturn, token.Identifier, token.Plus, token.Identifier,
		token.Semicolon, token.Newline,
		token.KwEnd,
		token.EOF,
	}

	got := make([]token.Kind, 0, len(tokens))
	for _, tok := range tokens {
		got = append(got, tok.Kind)
	}

	if len(got) != len(expectedKinds) {
		t.Fatalf("expected %d tokens, got %d\ngot:  %v\nwant: %v", len(expectedKinds), len(got), got, expectedKinds)
	}

	for i := range expectedKinds {
		if got[i] != expectedKinds[i] {
			t.Errorf("token[%d] = %s, want %s", i, got[i], expectedKinds[i])
		}
	}
}

func TestLexerStringWithNewlineErrorCorrected(t *testing.T) {
	// The lexString function currently emits wrong error "unterminated character literal"
	// for newlines inside strings. This test documents that behavior.
	input := "\"hello\nworld\""
	_, diag := lex(input)

	if !diag.HasErrors() {
		t.Error("expected error for newline in string literal")
	}
}
