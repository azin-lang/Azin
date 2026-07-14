package lexer

import "github.com/azin-lang/azin/compiler/internal/token"

func (l *Lexer) lexIdentifier(start token.Position) token.Token {
	l.consumeWhile(isIdentifierContinue)

	if kind, ok := token.Keywords[string(l.file.Slice(start.Offset, l.cursor))]; ok {
		return l.emit(kind, start)
	}

	return l.emit(token.Identifier, start)
}
