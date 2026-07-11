// Package diagnostics provides structures and engines to record, format,
// and report compiler errors, warnings, and informational notes.
package diagnostics

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/token"
)

// DiagnosticKind is the severity level of a diagnostic.
type DiagnosticKind uint8

const (
	Note = iota
	Warning
	Error
)

// A Diagnostic represents a single compiler message at a specific location.
type Diagnostic struct {
	Kind     DiagnosticKind //  kind
	Message  string         // diagnostic message
	Position token.Position // position at which the diagnostic was emited
	Length   uint32         // length of the offending
}

func (k DiagnosticKind) String() string {
	switch k {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Note:
		return "note"
	default:
		return "unknown"
	}
}

func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s: %s", d.Kind, d.Message)
}
