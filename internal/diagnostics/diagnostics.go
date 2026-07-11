package diagnostics

import (
	"fmt"
)

type DiagnosticKind uint8

const (
	Note = iota
	Warning
	Error
)

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

type Diagnostic struct {
	Kind    DiagnosticKind
	Message string
	Line    int
	Column  int
	Offset  int
	Length  int
}

func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s: %s", d.Kind, d.Message)
}
