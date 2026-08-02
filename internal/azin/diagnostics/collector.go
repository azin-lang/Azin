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

	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	yellowBold := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	for _, d := range c.List() {
		sevStr := d.Severity.String()
		styleColor := redBold
		pipeColor := cyanBold

		switch d.Severity {
		case SeverityError:
			sevStr = redBold("error")
			styleColor = redBold
			pipeColor = cyanBold
		case SeverityWarning:
			sevStr = yellowBold("warning")
			styleColor = yellowBold
			pipeColor = yellowBold
		case SeverityInfo:
			sevStr = cyanBold("info")
			styleColor = cyanBold
			pipeColor = cyanBold
		}

		codeStr := ""
		if d.Code != "" {
			codeStr = bold(fmt.Sprintf("[%s]", d.Code))
		}

		locStr := bold(d.Location.String())

		// Header
		if codeStr != "" {
			out.WriteString(fmt.Sprintf("%s: %s %s: %s\n", locStr, sevStr, codeStr, d.Message))
		} else {
			out.WriteString(fmt.Sprintf("%s: %s: %s\n", locStr, sevStr, d.Message))
		}

		// Render snippet with semantic colors
		if d.Location.File != nil && d.Location.Span.IsValidAndNonEmpty() {
			renderedLoc := source.PrintColoredDiagnostic(d.Location, d.Label, d.Note, d.Help, pipeColor, styleColor)
			out.WriteString(renderedLoc)
		}
		out.WriteString("\n")
	}

	return out.String()
}
