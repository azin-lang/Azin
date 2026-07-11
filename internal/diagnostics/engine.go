package diagnostics

import (
	"fmt"
	"strings"

	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

// Engine collects diagnostics for a source file.
type Engine struct {
	file        *source.File
	diagnostics []Diagnostic
	hasErrors   bool
}

// New returns a new Engine for the given file.
func New(file *source.File) *Engine {
	return &Engine{
		file: file,
	}
}

// Report adds a diagnostic to the engine.
func (e *Engine) Report(kind DiagnosticKind, pos token.Position, length uint32, format string, args ...any) {
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

// ReportError logs an error-level diagnostic.
func (e *Engine) ReportError(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Error, pos, length, format, args...)
}

// ReportWarning logs a warning-level diagnostic.
func (e *Engine) ReportWarning(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Warning, pos, length, format, args...)
}

// ReportNote logs a note-level diagnostic.
func (e *Engine) ReportNote(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Note, pos, length, format, args...)
}

// Diagnostics returns all recorded diagnostics.
func (e *Engine) Diagnostics() []Diagnostic {
	return e.diagnostics
}

// HasErrors reports whether any errors have been recorded.
func (e *Engine) HasErrors() bool {
	return e.hasErrors
}

// Err returns the engine as an error if it contains errors, or nil otherwise.
func (e *Engine) Err() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

// Formats all diagnostics into a readable string.
func (e *Engine) Error() string {
	var b strings.Builder

	for i, d := range e.diagnostics {
		if i > 0 {
			b.WriteString("\n\n")
		}

		line, column := e.file.LineColumn(d.Position.Offset)

		fmt.Fprintf(&b, "%s:%d:%d: %s: %s\n", e.file.Name(), line, column, d.Kind, d.Message)

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
