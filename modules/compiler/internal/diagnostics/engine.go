package diagnostics

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/azin-lang/azin/compiler/internal/source"
	"github.com/azin-lang/azin/compiler/internal/token"
)

// Engine collects diagnostics for a source file.
type Engine struct {
	mu          sync.RWMutex
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
	e.mu.Lock()
	defer e.mu.Unlock()

	if kind == Error {
		e.hasErrors = true
	}
	e.diagnostics = append(e.diagnostics, Diagnostic{
		Kind: kind, Message: fmt.Sprintf(format, args...), Position: pos, Length: length,
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
	e.mu.RLock()
	defer e.mu.RUnlock()

	res := make([]Diagnostic, len(e.diagnostics))
	copy(res, e.diagnostics)
	return res
}

// HasErrors reports whether any errors have been recorded.
func (e *Engine) HasErrors() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
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
	e.mu.RLock()
	defer e.mu.RUnlock()

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

		_, _ = fmt.Fprintf(&b, "%s:%d:%d: %s: %s\n", e.file.Name(), line, column, d.Kind, d.Message)

		_, _ = fmt.Fprintf(&b, "%*s |\n", gutter, "")
		_, _ = fmt.Fprintf(&b, "%*d | %s\n", gutter, line, src)
		_, _ = fmt.Fprintf(&b, "%*s | ", gutter, "")

		colIdx := int(column) - 1
		if colIdx < 0 {
			colIdx = 0
		}
		if colIdx > len(src) {
			colIdx = len(src)
		}

		prefix := src[:colIdx]
		for _, ch := range prefix {
			if ch == '\t' {
				b.WriteByte('\t')
			} else {
				b.WriteByte(' ')
			}
		}

		b.WriteByte('^')

		if d.Length > 0 {
			for i := uint32(1); i < d.Length; i++ {
				b.WriteByte('~')
			}
		}

		b.WriteByte('\n')
	}

	return b.String()
}
