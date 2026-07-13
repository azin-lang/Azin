// Package diagnostics provides structures and engines to record, format,
// and report compiler errors, warnings, and informational notes.
package diagnostics

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/token"
	"github.com/fatih/color"
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
		return color.New(color.FgRed, color.Bold).Sprint("error")
	case Warning:
		return color.New(color.FgHiBlue, color.Bold).Sprint("warning")
	case Note:
		return color.New(color.FgGreen, color.Bold).Sprint("note")
	default:
		return color.New(color.FgRed, color.Bold).Sprint("unknown")
	}
}

func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s: %s", d.Kind, d.Message)
}
