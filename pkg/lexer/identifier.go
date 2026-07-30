package lexer

import (
	token "github.com/azin-lang/Azin/pkg/token"
)

func (l *Lexer) lexIdentifier(start token.Position) token.Token {
	l.consumeWhile(isIdentifierContinue)

	identBytes := l.src[start.Offset:l.cursor]
	if kind, ok := token.Keywords[string(identBytes)]; ok {
		return l.emit(kind, start)
	}

	return l.emit(token.Identifier, start)
}
