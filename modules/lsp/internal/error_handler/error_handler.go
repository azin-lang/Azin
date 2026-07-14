package error_handler

import (
	"fmt"

	"github.com/azin-lang/azin/lsp/internal/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// LspErrorReporter implements parser.ErrorReporter
type LspErrorReporter struct {
	Diagnostics []protocol.Diagnostic
}

func (r *LspErrorReporter) ReportError(pos token.Position, length uint32, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)

	sourceName := "azin" // LSP positions are 0-indexed, while your token positions are likely 1-indexed
	startLine := pos.Offset - 1
	startChar := pos.Offset - 1

	endLine := startLine
	endChar := pos.Offset - 1 + length

	r.Diagnostics = append(r.Diagnostics, protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: startLine, Character: startChar},
			End:   protocol.Position{Line: endLine, Character: endChar},
		},
		Severity: new(protocol.DiagnosticSeverityError),
		Source:   &sourceName,
		Message:  msg,
	})
}

func (r *LspErrorReporter) Err() error {
	if len(r.Diagnostics) > 0 {
		return fmt.Errorf("encountered %d errors", len(r.Diagnostics))
	}
	return nil
}
