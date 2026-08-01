package diagnostics

import (
	"github.com/azin-lang/Azin/internal/azin/source"
)

// Diagnostic represents a single compiler error, warning, or notice.
type Diagnostic struct {
	// Code is an optional unique identifier for the diagnostic (e.g. "AZ0001").
	Code string

	// Location points to the exact file and span where the diagnostic applies.
	Location source.Location

	// Message is the descriptive explanation of the error or warning.
	Message string

	// Severity indicates whether the diagnostic is an error, warning, or info notice.
	Severity Severity
}
