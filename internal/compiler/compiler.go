package compiler

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/token"
)

func dumpTokens(tokens []token.Token, source []byte) {
	for _, token := range tokens {
		fmt.Println(token.Format(source))
	}
}

func runLexerStage(
	source []byte,
	printTokens bool,
	diag *diagnostics.Engine,
) []token.Token {
	lx := lexer.New(source, diag)
	tokens := lx.Tokenize()

	if printTokens {
		dumpTokens(tokens, source)
	}

	return tokens
}

func Compile(source []byte, printTokens bool) error {
	diag := diagnostics.New()
	_ = runLexerStage(source, printTokens, diag)

	/*
	 * parser := parser.New(lexer)
	 * ast, err := parser.Parse()
	 * if err != nil {
	 *     return err
	 * }
	 */

	return diag.Err()
}
