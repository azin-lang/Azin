//nolint:unparam
package diagnostics_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

func newTestEngine(text string) (*diagnostics.Engine, *source.File) {
	file := source.New("test.az", []byte(text))
	return diagnostics.New(file), file
}

func TestReportError(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	diag.ReportError(token.Position{Offset: 0}, 5, "test error")

	if !diag.HasErrors() {
		t.Error("HasErrors() = false after reporting error")
	}
	if diag.Err() == nil {
		t.Error("Err() = nil after reporting error")
	}
}

func TestReportWarning(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	diag.ReportWarning(token.Position{Offset: 0}, 5, "test warning")

	if diag.HasErrors() {
		t.Error("HasErrors() = true after reporting only warnings")
	}
}

func TestReportNote(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	diag.ReportNote(token.Position{Offset: 0}, 5, "test note")

	if diag.HasErrors() {
		t.Error("HasErrors() = true after reporting only notes")
	}
}

func TestNoErrors(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	if diag.HasErrors() {
		t.Error("HasErrors() = true with no diagnostics")
	}
	if diag.Err() != nil {
		t.Error("Err() != nil with no diagnostics")
	}
}

func TestDiagnosticsCollection(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	diag.ReportError(token.Position{Offset: 0}, 5, "first error")
	diag.ReportWarning(token.Position{Offset: 6}, 5, "a warning")
	diag.ReportError(token.Position{Offset: 0}, 3, "second error")

	all := diag.Diagnostics()
	if len(all) != 3 {
		t.Fatalf("got %d diagnostics, want 3", len(all))
	}

	if all[0].Kind != diagnostics.Error {
		t.Errorf("first diagnostic kind = %d, want %d", all[0].Kind, diagnostics.Error)
	}
	if all[1].Kind != diagnostics.Warning {
		t.Errorf("second diagnostic kind = %d, want %d", all[1].Kind, diagnostics.Warning)
	}
	if all[2].Kind != diagnostics.Error {
		t.Errorf("third diagnostic kind = %d, want %d", all[2].Kind, diagnostics.Error)
	}
}

func TestErrorOutputContainsSourceSnippet(t *testing.T) {
	diag, _ := newTestEngine("line one\nline two\nline three")
	diag.ReportError(token.Position{Offset: 0}, 4, "something went wrong")

	err := diag.Err()
	if err == nil {
		t.Fatal("expected error")
	}

	msg := err.Error()
	if !strings.Contains(msg, "test.az") {
		t.Errorf("error missing filename, got: %s", msg)
	}
	if !strings.Contains(msg, "line one") {
		t.Errorf("error missing source line, got: %s", msg)
	}
	if !strings.Contains(msg, "error") {
		t.Errorf("error missing severity, got: %s", msg)
	}
}

func TestConcurrentSafety(t *testing.T) {
	diag, _ := newTestEngine("hello world")

	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			diag.ReportError(token.Position{Offset: 0}, 5, "concurrent error")
			diag.HasErrors()
			diag.Diagnostics()
			_ = diag.Err()
		})
	}
	wg.Wait()

	all := diag.Diagnostics()
	if len(all) != 20 {
		t.Errorf("got %d diagnostics, want 20", len(all))
	}
}

func TestErrorOutputFormatsPosition(t *testing.T) {
	diag, _ := newTestEngine("fn main() do\n    return 0;\nend")
	diag.ReportError(token.Position{Offset: 3}, 4, "something wrong")

	msg := diag.Err().Error()
	if !strings.Contains(msg, "test.az:1:4") {
		t.Errorf("error missing correct position, got: %s", msg)
	}
}

func TestMultipleErrors(t *testing.T) {
	diag, _ := newTestEngine("line1\nline2\nline3")
	diag.ReportError(token.Position{Offset: 0}, 3, "first")
	diag.ReportError(token.Position{Offset: 6}, 3, "second")

	msg := diag.Err().Error()
	lines := strings.Count(msg, "error:")
	if lines != 2 {
		t.Errorf("expected 2 errors in output, got %d", lines)
	}
}

func TestWarningWithoutError(t *testing.T) {
	diag, _ := newTestEngine("hello")
	diag.ReportWarning(token.Position{Offset: 0}, 5, "just a warning")

	if diag.HasErrors() {
		t.Error("HasErrors() should be false with only warnings")
	}
	if diag.Err() != nil {
		t.Error("Err() should be nil with only warnings")
	}
}

func TestErrorLimitSuppressesExcess(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	diag.SetErrorLimit(3)

	for i := range 10 {
		diag.ReportError(token.Position{Offset: 0}, 5, "error %d", i)
	}

	all := diag.Diagnostics()
	if len(all) > 4 {
		t.Errorf("expected at most 4 diagnostics (3 errors + 1 limit note), got %d", len(all))
	}

	limitNote := false
	for _, d := range all {
		if d.Kind == diagnostics.Note && strings.Contains(d.Message, "too many errors") {
			limitNote = true
			break
		}
	}
	if !limitNote {
		t.Error("expected 'too many errors' limit note")
	}
}

func TestErrorLimitUnlimited(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	diag.SetErrorLimit(0)

	for i := range 100 {
		diag.ReportError(token.Position{Offset: 0}, 5, "error %d", i)
	}

	all := diag.Diagnostics()
	if len(all) != 100 {
		t.Errorf("expected 100 diagnostics, got %d", len(all))
	}
}

func TestErrorLimitWarningsNotAffected(t *testing.T) {
	diag, _ := newTestEngine("hello world")
	diag.SetErrorLimit(2)

	for i := range 5 {
		diag.ReportError(token.Position{Offset: 0}, 5, "error %d", i)
	}
	for i := range 5 {
		diag.ReportWarning(token.Position{Offset: 0}, 5, "warning %d", i)
	}

	all := diag.Diagnostics()
	warnings := 0
	for _, d := range all {
		if d.Kind == diagnostics.Warning {
			warnings++
		}
	}
	if warnings != 5 {
		t.Errorf("expected 5 warnings, got %d", warnings)
	}
}

func TestErrorLimitDefault(t *testing.T) {
	diag, _ := newTestEngine("hello world")

	// Default limit is 50, so 60 errors should trigger the limit
	for i := range 60 {
		diag.ReportError(token.Position{Offset: 0}, 5, "error %d", i)
	}

	all := diag.Diagnostics()
	// 50 errors + 1 limit note = 51 max
	if len(all) > 51 {
		t.Errorf("expected at most 51 diagnostics, got %d", len(all))
	}
}
