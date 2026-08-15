package lexer

import (
	"fmt"
	"unicode"

	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/text"
)

func (l *Lexer) scanSyntaxToken(startOffset int) syntax.Kind {
	ch, _ := l.reader.Next()

	switch ch {
	case '(':
		return syntax.OpenParenToken
	case ')':
		return syntax.CloseParenToken
	case '{':
		return syntax.OpenBraceToken
	case '}':
		return syntax.CloseBraceToken
	case '[':
		return syntax.OpenBracketToken
	case ']':
		return syntax.CloseBracketToken
	case ',':
		return syntax.CommaToken
	case ':':
		return syntax.ColonToken
	case ';':
		return syntax.SemicolonToken
	case '.':
		return syntax.DotToken
	case '"':
		return l.scanString(startOffset)

	case '+':
		if l.match('=') {
			return syntax.PlusEqualsToken
		}
		return syntax.PlusToken
	case '-':
		if l.match('=') {
			return syntax.MinusEqualsToken
		}
		return syntax.MinusToken
	case '*':
		if l.match('=') {
			return syntax.StarEqualsToken
		}
		return syntax.StarToken
	case '/':
		if l.match('=') {
			return syntax.SlashEqualsToken
		}
		return syntax.SlashToken
	case '=':
		if l.match('=') {
			return syntax.EqualsEqualsToken
		}
		return syntax.EqualsToken
	case '!':
		if l.match('=') {
			return syntax.BangEqualsToken
		}
		return syntax.BangToken
	case '<':
		if l.match('=') {
			return syntax.LessEqualsToken
		}
		return syntax.LessToken
	case '>':
		if l.match('=') {
			return syntax.GreaterEqualsToken
		}
		return syntax.GreaterToken
	case '&':
		if l.match('&') {
			return syntax.AmpersandAmpersandToken
		}
		return syntax.AmpersandToken

	case '|':
		if l.match('|') {
			return syntax.PipePipeToken
		}
		return syntax.PipeToken
	}

	if unicode.IsLetter(ch) || ch == '_' {
		// The first character is already consumed; this will consume the rest.
		return l.scanIdentifierOrKeyword(startOffset)
	}

	if unicode.IsDigit(ch) {
		// We already consumed the first digit, but scanNumber expects to read it.
		// Teleport back to exactly where this token started.
		l.reader.Restore(startOffset)
		return l.scanNumber(startOffset)
	}

	l.diagnostics.ErrorWithNoteAndHelp(
		l.loc(startOffset, l.reader.Offset()),
		"E001",
		fmt.Sprintf("I don't know what to do with '%c'.", ch),
		"This character isn't part of any valid Azin syntax that I know of.",
		"Maybe it's a typo? Try removing it or replacing it with a valid operator.",
		"Unexpected character",
	)
	return syntax.BadToken
}

func (l *Lexer) consumeAlphanumeric() {
	for !l.reader.EOF() {
		ch, _ := l.reader.Peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			l.reader.Next()
		} else {
			break
		}
	}
}

func (l *Lexer) consumeDigits() {
	for !l.reader.EOF() {
		ch, _ := l.reader.Peek()
		if unicode.IsDigit(ch) {
			l.reader.Next()
		} else {
			break
		}
	}
}

func (l *Lexer) scanIdentifierOrKeyword(startOffset int) syntax.Kind {
	l.consumeAlphanumeric()
	source := l.file.Text(text.SpanFromBounds(startOffset, l.reader.Offset()))
	return syntax.LookupKeyword(source)
}

func (l *Lexer) scanNumber(startOffset int) syntax.Kind {
	isFloat := false
	l.consumeDigits()

	ch, _ := l.reader.Peek()
	if ch == '.' {
		// Checkpoint before eating the dot
		dotCheckpoint := l.reader.Checkpoint()

		l.reader.Next()
		nextCh, _ := l.reader.Peek()

		if unicode.IsDigit(nextCh) {
			isFloat = true
			l.consumeDigits()
		} else {
			// False alarm (e.g. `1.method()`). Undo consuming the dot.
			l.reader.Restore(dotCheckpoint)
		}
	}

	ch, _ = l.reader.Peek()
	if ch == 'e' || ch == 'E' {
		// Checkpoint before eating 'e', because we might also eat a '+' or '-'
		// and still fail to find digits.
		expCheckpoint := l.reader.Checkpoint()

		l.reader.Next() // Consume 'e'

		hasSign := false
		signCh, _ := l.reader.Peek()
		if signCh == '+' || signCh == '-' {
			l.reader.Next() // Consume sign
			hasSign = true
		}

		expCh, _ := l.reader.Peek()
		if unicode.IsDigit(expCh) {
			isFloat = true
			l.consumeDigits()
		} else {
			if hasSign {
				l.consumeAlphanumeric()
				l.diagnostics.ErrorWithNoteAndHelp(
					l.loc(startOffset, l.reader.Offset()),
					"E006",
					"Missing digits after exponent sign.",
					"Scientific notation needs numbers after the 'e+' or 'e-'.",
					"Try adding the exponent value (like '1e+5').",
					"Invalid scientific notation",
				)
				return syntax.BadToken
			}
			// It was just an 'e' (or 'E') attached to a number, but not scientific notation.
			// Restore the reader to exactly before we ate the 'e'.
			l.reader.Restore(expCheckpoint)
		}
	}

	ch, _ = l.reader.Peek()
	if unicode.IsLetter(ch) || ch == '_' {
		l.consumeAlphanumeric()
		l.diagnostics.ErrorWithNoteAndHelp(
			l.loc(startOffset, l.reader.Offset()),
			"E004",
			"This looks like a number mixed with a word.",
			"I get confused when letters come right after numbers without a space.",
			"If this is a number followed by a variable, try adding a space between them (like '18 years').",
			"Number format issue",
		)
		return syntax.BadToken
	}

	if isFloat {
		return syntax.FloatLiteralToken
	}
	return syntax.IntegerLiteralToken
}

func (l *Lexer) scanString(startOffset int) syntax.Kind {
	for !l.reader.EOF() {
		ch, _ := l.reader.Next()

		if ch == '\\' && !l.reader.EOF() {
			escapeStart := l.reader.Offset() - 1
			esc, _ := l.reader.Next()

			if !isValidEscape(esc) {
				l.diagnostics.ErrorWithNoteAndHelp(
					l.loc(escapeStart, l.reader.Offset()),
					"E005",
					fmt.Sprintf("I don't recognize the escape sequence '\\%c'.", esc),
					"It looks like you're trying to escape a character, but I only know a specific set.",
					"The escape sequences I understand are \\n, \\r, \\t, \\\\, \\\", and \\0.",
					"Unknown escape sequence",
				)
			}
			continue
		}

		if ch == '"' {
			return syntax.StringLiteralToken
		}
		if ch == '\n' || ch == '\r' {
			break
		}
	}

	l.diagnostics.ErrorWithNoteAndHelp(
		l.loc(startOffset, l.reader.Offset()),
		"E002",
		"This string is missing a closing quote!",
		"I hit the end of the line while still reading your string. Strings can't jump to a new line on their own.",
		"Try adding a '\"' at the end of this line!",
		"Unfinished string",
	)

	return syntax.StringLiteralToken
}

func isValidEscape(ch rune) bool {
	switch ch {
	case 'n', 'r', 't', '\\', '"', '0':
		return true
	}
	return false
}
