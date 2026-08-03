package main

import (
	"fmt"
	"strings"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/lexer"
	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
)

func main() {
	code := []byte(`importc "stdio.h"

enum Types is
   Teacher
   Doctor
   None
end

struct Human is
   mut age: int
   mut name: string
   mut type: Types
   @  // E001: Unexpected character ('@')
end
`)

	file := source.NewSourceText("test.az", 1, code)
	diags := diagnostics.NewCollector()
	lex := lexer.New(file, diags)

	tokenIndex := 0
	for {
		tok := lex.NextToken()

		leadingText := ""
		if tok.LeadingTrivia() != nil {
			// Type-assert to extract text safely from the green.Node interface
			if tr, ok := tok.LeadingTrivia().(*green.Trivia); ok {
				leadingText = tr.Text()
			}
		}

		trailingText := ""
		if tok.TrailingTrivia() != nil {
			if tr, ok := tok.TrailingTrivia().(*green.Trivia); ok {
				trailingText = tr.Text()
			}
		}

		fmt.Printf("[%03d] %v\n", tokenIndex, tok.Kind())

		fmt.Printf("      ├── Leading trivia:  %s\n", formatMultiLineTrivia(leadingText))
		fmt.Printf("      ├── Token text:      %q\n", tok.Text())
		fmt.Printf("      └── Trailing trivia: %s\n", formatMultiLineTrivia(trailingText))
		fmt.Println()

		if tok.Kind() == syntax.EndOfFileToken {
			break
		}
		tokenIndex++
	}

	fmt.Println(strings.Repeat("═", 66))

	if diags.HasErrors() {
		fmt.Println(diags.PrintAll())
	} else {
		fmt.Println("Compilation successful with zero errors.")
	}
}

// formatMultiLineTrivia formats newlines neatly when displayed in an indented tree format.
func formatMultiLineTrivia(s string) string {
	if s == "" {
		return "∅"
	}
	// If the trivia has newlines, format them cleanly with proper indentation
	replacer := strings.NewReplacer("\r\n", "\\n", "\n", "\\n", "\r", "\\r", "\t", "\\t")
	return fmt.Sprintf("%q", replacer.Replace(s))
}
