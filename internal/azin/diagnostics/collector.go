package diagnostics

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/fatih/color"
	"golang.org/x/term"
)

// Collector thread-safely gathers diagnostics throughout compilation passes.
type Collector struct {
	mu          sync.Mutex
	diagnostics []Diagnostic
	hasErrors   bool
}

// NewCollector creates an empty diagnostic collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Error records a new compilation error with an optional error code.
func (c *Collector) Error(loc source.Location, code, format string, args ...any) {
	c.add(loc, code, fmt.Sprintf(format, args...), SeverityError, "", "")
}

// ErrorWithHelp records a compilation error accompanied by an inline label and an actionable help hint.
func (c *Collector) ErrorWithHelp(loc source.Location, code, label, help, format string, args ...any) {
	c.add(loc, code, fmt.Sprintf(format, args...), SeverityError, label, help)
}

// Warn records a new compilation warning with an optional error code.
func (c *Collector) Warn(loc source.Location, code, format string, args ...any) {
	c.add(loc, code, fmt.Sprintf(format, args...), SeverityWarning, "", "")
}

// Info records a new compilation information message with an optional code.
func (c *Collector) Info(loc source.Location, code, format string, args ...any) {
	c.add(loc, code, fmt.Sprintf(format, args...), SeverityInfo, "", "")
}

func (c *Collector) add(loc source.Location, code, msg string, sev Severity, label, help string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if sev == SeverityError {
		c.hasErrors = true
	}

	c.diagnostics = append(c.diagnostics, Diagnostic{
		Code:     code,
		Location: loc,
		Message:  msg,
		Severity: sev,
		Label:    label,
		Help:     help,
	})
}

// HasErrors reports whether any error-level diagnostics have been logged.
func (c *Collector) HasErrors() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hasErrors
}

// Clear resets the collector state, wiping all accumulated diagnostics.
func (c *Collector) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.diagnostics = nil
	c.hasErrors = false
}

// List returns a sorted, cloned copy of all collected diagnostics.
func (c *Collector) List() []Diagnostic {
	c.mu.Lock()
	defer c.mu.Unlock()

	list := slices.Clone(c.diagnostics)
	slices.SortFunc(list, func(a, b Diagnostic) int {
		pathA := ""
		if a.Location.File != nil {
			pathA = a.Location.File.Path()
		}
		pathB := ""
		if b.Location.File != nil {
			pathB = b.Location.File.Path()
		}

		if pathA != pathB {
			return cmp.Compare(pathA, pathB)
		}
		return cmp.Compare(a.Location.Span.Start, b.Location.Span.Start)
	})

	return list
}

// ErrorWithNoteAndHelp records an error with an inline label, an actionable note, and a help hint.
func (c *Collector) ErrorWithNoteAndHelp(loc source.Location, code, label, note, help, format string, args ...any) {
	c.addFull(loc, code, fmt.Sprintf(format, args...), SeverityError, label, note, help)
}

func (c *Collector) addFull(loc source.Location, code, msg string, sev Severity, label, note, help string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if sev == SeverityError {
		c.hasErrors = true
	}

	c.diagnostics = append(c.diagnostics, Diagnostic{
		Code:     code,
		Location: loc,
		Message:  msg,
		Severity: sev,
		Label:    label,
		Note:     note,
		Help:     help,
	})
}

// PrintAll renders all diagnostics with fully color-coordinated terminal styling.
func (c *Collector) PrintAll() string {
	var out strings.Builder

	if _, present := os.LookupEnv("NO_COLOR"); present {
		color.NoColor = true
	}

	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).Add(color.Bold).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	for i, d := range c.List() {
		if i > 0 {
			out.WriteString("\n\n") // Space between multiple errors
		}

		styleColor := cyan
		pipeColor := dim

		switch d.Severity {
		case SeverityError:
			styleColor = red
		case SeverityWarning:
			styleColor = yellow
		case SeverityInfo:
			styleColor = cyan
		}

		title := strings.ToUpper(d.Message)

		termWidth := 80
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			termWidth = width
		}

		maxWidth := 100
		if termWidth > maxWidth {
			termWidth = maxWidth
		}

		dashCount := max(termWidth-len(title)-4, 2)

		banner := fmt.Sprintf("%s %s %s", cyan("──"), cyan(title), cyan(strings.Repeat("─", dashCount)))
		out.WriteString(banner + "\n\n")

		if d.Location.File != nil && d.Location.Span.IsValidAndNonEmpty() {
			renderedLoc := source.PrintColoredDiagnostic(
				d.Location, d.Label, d.Note, d.Help,
				termWidth,
				pipeColor, styleColor, bold,
			)
			out.WriteString(renderedLoc)
		}
	}

	return out.String()
}
