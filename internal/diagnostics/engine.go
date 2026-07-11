package diagnostics

import (
	"fmt"
	"strings"

	"github.com/azin-lang/Azin/internal/source" //
	"github.com/azin-lang/Azin/internal/token"
)

// the Engine struct is responsible for reporting diagnostics.
type Engine struct {
	file        *source.File
	diagnostics []Diagnostic
	hasErrors   bool
}

// creates a new diagnostic engine for the given file.
func New(file *source.File) *Engine {
	return &Engine{
		file: file,
	}
}

// reports a diagnostic message.
func (e *Engine) Report(
	kind DiagnosticKind,
	pos token.Position,
	length uint32,
	format string,
	args ...any,
) {
	if kind == Error {
		e.hasErrors = true
	}

	e.diagnostics = append(e.diagnostics, Diagnostic{
		Kind:     kind,
		Message:  fmt.Sprintf(format, args...),
		Position: pos,
		Length:   length,
	})
}

// reports an error diagnostic message.
func (e *Engine) ReportError(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Error, pos, length, format, args...)
}

// reports a warning diagnostic message.
func (e *Engine) ReportWarning(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Warning, pos, length, format, args...)
}

// reports a note diagnostic message.
func (e *Engine) ReportNote(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Note, pos, length, format, args...)
}

// returns the list of diagnostics.
func (e *Engine) Diagnostics() []Diagnostic {
	return e.diagnostics
}

// returns whether the engine has errors.
func (e *Engine) HasErrors() bool {
	return e.hasErrors
}

// returns the error if the engine has errors, otherwise nil.
func (e *Engine) Err() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

// returns the line and column of the given position.
func (e *Engine) LineColumn(pos token.Position) (line, column uint32) {
	return e.file.LineColumn(pos.Offset)
}

// returns the text of the given token.
func (e *Engine) Text(tok token.Token) []byte {
	return e.file.Text(tok)
}

// returns the line of the given line number.
func (e *Engine) Line(line uint32) []byte {
	return e.file.Line(line)
}

// returns the error message if the engine has errors, otherwise nil.
func (e *Engine) Error() string {
	var b strings.Builder

	for i, d := range e.diagnostics {
		if i > 0 {
			b.WriteString("\n\n")
		}

		line, column := e.file.LineColumn(d.Position.Offset)

		fmt.Fprintf(
			&b,
			"%s:%d:%d: %s: %s\n",
			e.file.Name(),
			line,
			column,
			d.Kind,
			d.Message,
		)

		src := e.file.Line(line)
		b.Write(src)
		b.WriteByte('\n')

		prefix := src[:column-1]

		for _, ch := range prefix {
			if ch == '\t' {
				b.WriteByte('\t')
			} else {
				b.WriteByte(' ')
			}
		}

		b.WriteByte('^')
		for i := uint32(1); i < d.Length; i++ {
			b.WriteByte('~')
		}
	}

	return b.String()
}
