package main

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/lexer"
	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

func main() {
	sourceCode := `
	// Compute some math
	var      x= 42 @#$;fn compute()do
// return result
		return x;
	end
	`

	file := source.New("test.az", 1, []byte(sourceCode))
	collector := diagnostics.NewCollector()

	lex := lexer.New(file, collector)

	fmt.Println("==================================================")
	fmt.Println("            AZIN LEXER INSPECTOR                  ")
	fmt.Println("==================================================")

	for {
		tok := lex.NextToken()

		for _, t := range tok.LeadingTrivia {
			fmt.Printf("  [Leading]  %-24s | Text: %q\n", kindName(t.Kind), file.Text(t.Span))
		}

		tokenText := tok.Text(file)
		fmt.Printf("Token ->     %-24s | Text: %-12q | Span: %v\n", kindName(tok.Kind), tokenText, tok.Span)

		for _, t := range tok.TrailingTrivia {
			fmt.Printf("  [Trailing] %-24s | Text: %q\n", kindName(t.Kind), file.Text(t.Span))
		}

		if tok.Kind == syntax.EndOfFile {
			break
		}
	}

	fmt.Println("==================================================")

	if collector.HasErrors() {
		fmt.Println("\nCompilation diagnostics:\n")
		fmt.Print(collector.PrintAll())
	} else {
		fmt.Println("Lexing completed successfully with zero errors.")
	}
}

// kindName maps raw syntax.Kind integers to human-readable names.
func kindName(k syntax.Kind) string {
	switch k {
	case syntax.Unknown:
		return "Unknown"
	case syntax.EndOfFile:
		return "EndOfFile"
	case syntax.Error:
		return "Error"
	case syntax.WhitespaceTrivia:
		return "WhitespaceTrivia"
	case syntax.EndOfLineTrivia:
		return "EndOfLineTrivia"
	case syntax.SingleLineCommentTrivia:
		return "SingleLineCommentTrivia"
	case syntax.IdentifierToken:
		return "IdentifierToken"
	case syntax.IntegerLiteralToken:
		return "IntegerLiteralToken"
	case syntax.FloatLiteralToken:
		return "FloatLiteralToken"
	case syntax.CharacterLiteralToken:
		return "CharacterLiteralToken"
	case syntax.BooleanLiteralToken:
		return "BooleanLiteralToken"
	case syntax.StringLiteralToken:
		return "StringLiteralToken"
	case syntax.KeywordDefer:
		return "KeywordDefer"
	case syntax.KeywordDo:
		return "KeywordDo"
	case syntax.KeywordElse:
		return "KeywordElse"
	case syntax.KeywordEnd:
		return "KeywordEnd"
	case syntax.KeywordEnum:
		return "KeywordEnum"
	case syntax.KeywordFn:
		return "KeywordFn"
	case syntax.KeywordFor:
		return "KeywordFor"
	case syntax.KeywordIf:
		return "KeywordIf"
	case syntax.KeywordImport:
		return "KeywordImport"
	case syntax.KeywordImportC:
		return "KeywordImportC"
	case syntax.KeywordIs:
		return "KeywordIs"
	case syntax.KeywordLoop:
		return "KeywordLoop"
	case syntax.KeywordMut:
		return "KeywordMut"
	case syntax.KeywordReturn:
		return "KeywordReturn"
	case syntax.KeywordStop:
		return "KeywordStop"
	case syntax.KeywordStruct:
		return "KeywordStruct"
	case syntax.KeywordThen:
		return "KeywordThen"
	case syntax.KeywordType:
		return "KeywordType"
	case syntax.KeywordVar:
		return "KeywordVar"
	case syntax.KeywordWhile:
		return "KeywordWhile"
	case syntax.PlusToken:
		return "PlusToken"
	case syntax.PlusEqualsToken:
		return "PlusEqualsToken"
	case syntax.MinusToken:
		return "MinusToken"
	case syntax.MinusEqualsToken:
		return "MinusEqualsToken"
	case syntax.EqualsToken:
		return "EqualsToken"
	case syntax.EqualsEqualsToken:
		return "EqualsEqualsToken"
	case syntax.StarToken:
		return "StarToken"
	case syntax.StarEqualsToken:
		return "StarEqualsToken"
	case syntax.OpenParenToken:
		return "OpenParenToken"
	case syntax.CloseParenToken:
		return "CloseParenToken"
	case syntax.OpenBracketToken:
		return "OpenBracketToken"
	case syntax.CloseBracketToken:
		return "CloseBracketToken"
	case syntax.OpenBraceToken:
		return "OpenBraceToken"
	case syntax.CloseBraceToken:
		return "CloseBraceToken"
	case syntax.CommaToken:
		return "CommaToken"
	case syntax.SemicolonToken:
		return "SemicolonToken"
	case syntax.ColonToken:
		return "ColonToken"
	case syntax.DotToken:
		return "DotToken"
	case syntax.NewlineToken:
		return "NewlineToken"
	default:
		return fmt.Sprintf("Kind(%d)", k)
	}
}
