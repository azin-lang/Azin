package diagnostics

import (
	"github.com/azin-lang/Azin/internal/azin/text"
)

// Diagnostic represents a single compiler error, warning, or notice.
type Diagnostic struct {
	Code     string
	Location text.Location
	Message  string
	Severity Severity
	Label    string // Inline annotation next to the underline
	Note     string // Actionable suggestion (e.g. "try replacing with 'end'")
	Help     string // Explanatory guidance
}
