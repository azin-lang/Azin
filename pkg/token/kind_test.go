package token_test

import (
	"testing"

	tok "github.com/azin-lang/Azin/pkg/token"
)

func TestKindDisplayName(t *testing.T) {
	tests := []struct {
		kind tok.Kind
		want string
	}{
		{tok.Unknown, "unknown"},
		{tok.Identifier, "identifier"},
		{tok.IntegerLiteral, "integer literal"},
		{tok.FloatLiteral, "float literal"},
		{tok.StringLiteral, "string literal"},
		{tok.CharacterLiteral, "character literal"},
		{tok.KwFn, "'fn'"},
		{tok.KwDo, "'do'"},
		{tok.KwVar, "'var'"},
		{tok.KwMut, "'mut'"},
		{tok.KwReturn, "'return'"},
		{tok.KwEnd, "'end'"},
		{tok.KwIf, "'if'"},
		{tok.KwThen, "'then'"},
		{tok.KwElse, "'else'"},
		{tok.KwStruct, "'struct'"},
		{tok.KwIs, "'is'"},
		{tok.KwImportC, "'importC'"},
		{tok.KwImport, "'import'"},
		{tok.KwWhile, "'while'"},
		{tok.KwChar, "'char'"},
		{tok.KwInt, "'int'"},
		{tok.KwBool, "'bool'"},
		{tok.KwNull, "'null'"},
		{tok.KwUnit, "'unit'"},
		{tok.KwString, "'string'"},
		{tok.KwFloat, "'float'"},
		{tok.Plus, "'+'"},
		{tok.Minus, "'-'"},
		{tok.Star, "'*'"},
		{tok.Slash, "'/'"},
		{tok.Equal, "'='"},
		{tok.EqualEqual, "'=='"},
		{tok.Bang, "'!'"},
		{tok.BangEqual, "'!='"},
		{tok.Less, "'<'"},
		{tok.LessEqual, "'<='"},
		{tok.Greater, "'>'"},
		{tok.GreaterEqual, "'>='"},
		{tok.LeftParen, "'('"},
		{tok.RightParen, "')'"},
		{tok.Comma, "','"},
		{tok.Colon, "':'"},
		{tok.Semicolon, "';'"},
		{tok.Dot, "'.'"},
		{tok.Newline, "newline"},
		{tok.EOF, "end of file"},
	}

	for _, tt := range tests {
		got := tt.kind.DisplayName()
		if got != tt.want {
			t.Errorf("DisplayName(%d) = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestKindDisplayNameNonEmpty(t *testing.T) {
	for k := tok.Unknown; k <= tok.Error; k++ {
		name := k.DisplayName()
		if name == "" {
			t.Errorf("DisplayName(%d) is empty", k)
		}
	}
}

func TestKindHasText(t *testing.T) {
	tests := []struct {
		kind tok.Kind
		want bool
	}{
		{tok.Identifier, true},
		{tok.IntegerLiteral, true},
		{tok.FloatLiteral, true},
		{tok.StringLiteral, true},
		{tok.CharacterLiteral, true},
		{tok.Plus, false},
		{tok.KwFn, false},
		{tok.EOF, false},
	}

	for _, tt := range tests {
		got := tt.kind.HasText()
		if got != tt.want {
			t.Errorf("HasText(%d) = %v, want %v", tt.kind, got, tt.want)
		}
	}
}
