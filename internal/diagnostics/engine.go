package diagnostics

import (
	"fmt"
	"strings"
)

type Engine struct {
	diagnostics []Diagnostic
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) Report(kind DiagnosticKind, format string, args ...any) {
	e.diagnostics = append(e.diagnostics, Diagnostic{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	})
}

func (e *Engine) ReportError(format string, args ...any) {
	e.Report(Error, format, args...)
}

func (e *Engine) ReportWarning(format string, args ...any) {
	e.Report(Warning, format, args...)
}

func (e *Engine) ReportNote(format string, args ...any) {
	e.Report(Note, format, args...)
}

func (e *Engine) Diagnostics() []Diagnostic {
	return e.diagnostics
}

func (e *Engine) HasErrors() bool {
	for _, d := range e.diagnostics {
		if d.Kind == Error {
			return true
		}
	}

	return false
}
func (e *Engine) Err() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

func (e *Engine) Error() string {
	if len(e.diagnostics) == 0 {
		return ""
	}

	var b strings.Builder
	for i, d := range e.diagnostics {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(d.Error())
	}

	return b.String()
}
