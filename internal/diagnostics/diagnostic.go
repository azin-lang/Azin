package diagnostics

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/token"
)

// represents the kind of diagnostic.
type DiagnosticKind uint8

const (
	Note = iota
	Warning
	Error
)

// represents a diagnostic message.
type Diagnostic struct {
	Kind     DiagnosticKind
	Message  string
	Position token.Position
	Length   uint32
}

// returns the string representation of the diagnostic kind.
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

// returns the string representation of the diagnostic.
func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s: %s", d.Kind, d.Message)
}
