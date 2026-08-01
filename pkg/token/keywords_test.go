package token_test

import (
	"testing"

	tok "github.com/azin-lang/Azin/pkg/token"
)

func TestKeywordsContainAllRegistered(t *testing.T) {
	expected := map[string]tok.Kind{
		"fn":      tok.KwFn,
		"do":      tok.KwDo,
		"var":     tok.KwVar,
		"mut":     tok.KwMut,
		"return":  tok.KwReturn,
		"end":     tok.KwEnd,
		"char":    tok.KwChar,
		"int":     tok.KwInt,
		"bool":    tok.KwBool,
		"unit":    tok.KwUnit,
		"string":  tok.KwString,
		"float":   tok.KwFloat,
		"if":      tok.KwIf,
		"then":    tok.KwThen,
		"else":    tok.KwElse,
		"struct":  tok.KwStruct,
		"is":      tok.KwIs,
		"import":  tok.KwImport,
		"importc": tok.KwImportC,
		"loop":    tok.KwLoop,
		"while":   tok.KwWhile,
		"stop":    tok.KwStop,
		"null":    tok.KwNull,
		"enum":    tok.KwEnum,
		"defer":   tok.KwDefer,
	}

	for word, kind := range expected {
		got, ok := tok.Keywords[word]
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
		"is": true, "import": true, "importc": true, "loop": true, "while": true, "stop": true,
		"null": true, "enum": true, "defer": true,
	}

	for word := range tok.Keywords {
		if !known[word] {
			t.Errorf("Unexpected keyword entry: %q", word)
		}
	}

	if len(tok.Keywords) != len(known) {
		t.Errorf("Keywords map has %d entries, want %d", len(tok.Keywords), len(known))
	}
}
