package diagnostics

import (
	"fmt"
	"strconv"
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
	if len(e.diagnostics) == 0 {
		return ""
	}

	var b strings.Builder

	maxLine := 1

	for _, d := range e.diagnostics {
		line, _ := e.file.LineColumn(d.Position.Offset)
		if int(line) > maxLine {
			maxLine = int(line)
		}
	}

	gutter := len(strconv.Itoa(maxLine))

	for i, d := range e.diagnostics {
		if i > 0 {
			b.WriteByte('\n')
		}
		line, column := e.file.LineColumn(d.Position.Offset)
		src := e.file.Line(line)

		_, _ = fmt.Fprintf(&b, "%s: %s\n", d.Kind, d.Message)
		_, _ = fmt.Fprintf(&b, " --> %s:%d:%d\n", e.file.Name(), line, column)

		_, _ = fmt.Fprintf(&b, "%*s |\n", gutter, "")
		_, _ = fmt.Fprintf(&b, "%*d | %s\n", gutter, line, src)
		_, _ = fmt.Fprintf(&b, "%*s | ", gutter, "")

		prefix := src[:column-1]
		for _, ch := range prefix {
			if ch == '\t' {
				b.WriteByte('\t')
			} else {
				b.WriteByte(' ')
			}
		}

		if d.Length == 0 {
			b.WriteByte('^')
		} else {
			b.WriteByte('^')
			for i := uint32(1); i < d.Length; i++ {
				b.WriteByte('~')
			}
		}

		b.WriteByte(' ')
		b.WriteString(d.Message)
		b.WriteByte('\n')
	}

	return b.String()
}
