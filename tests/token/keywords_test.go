package token_test

import (
	"testing"

	"github.com/azin-lang/Azin/internal/token"
)

func TestKeywordsContainAllRegistered(t *testing.T) {
	expected := map[string]token.Kind{
		"fn":      token.KwFn,
		"do":      token.KwDo,
		"var":     token.KwVar,
		"mut":     token.KwMut,
		"return":  token.KwReturn,
		"end":     token.KwEnd,
		"char":    token.KwChar,
		"int":     token.KwInt,
		"bool":    token.KwBool,
		"unit":    token.KwUnit,
		"string":  token.KwString,
		"float":   token.KwFloat,
		"if":      token.KwIf,
		"then":    token.KwThen,
		"else":    token.KwElse,
		"struct":  token.KwStruct,
		"is":      token.KwIs,
		"importc": token.KwImportC,
		"loop":    token.KwLoop,
		"stop":    token.KwStop,
		"null":    token.KwNull,
		"enum":    token.KwEnum,
		"defer":   token.KwDefer,
	}

	for word, kind := range expected {
		got, ok := token.Keywords[word]
		if !ok {
			t.Errorf("Keywords map missing entry for %q", word)
			continue
		}
		if got != kind {
			t.Errorf("Keywords[%q] = %d, want %d", word, got, kind)
		}
	}
}

func TestKeywordsNoExtraEntries(t *testing.T) {
	known := map[string]bool{
		"fn": true, "do": true, "var": true, "mut": true,
		"return": true, "end": true, "char": true, "int": true,
		"bool": true, "unit": true, "string": true, "float": true,
		"if": true, "then": true, "else": true, "struct": true,
		"is": true, "importc": true, "loop": true, "stop": true,
		"null": true, "enum": true, "defer": true,
	}

	for word := range token.Keywords {
		if !known[word] {
			t.Errorf("Unexpected keyword entry: %q", word)
		}
	}

	if len(token.Keywords) != len(known) {
		t.Errorf("Keywords map has %d entries, want %d", len(token.Keywords), len(known))
	}
}
