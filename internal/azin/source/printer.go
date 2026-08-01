package source

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// PrintDiagnostic generates a formatted, CLI-friendly visual representation of a location.
// It extracts the surrounding line of code and uses caret symbols (^) to point
// out the exact span of characters associated with an error or warning.
func PrintDiagnostic(loc Location) string {
	if loc.File == nil || !loc.Span.IsValidAndNonEmpty() {
		return "<invalid location>"
	}

	pos := loc.Position()
	line := loc.File.Line(pos.Line)
	lineText := loc.File.BytesOf(NewSpan(line.Start, line.End))

	maxSquiggly := max(len(lineText)-int(pos.ByteColumn-1), 1)
	squigglyLen := int(loc.Span.Len())
	if squigglyLen > maxSquiggly {
		squigglyLen = maxSquiggly
	} else if squigglyLen <= 0 {
		squigglyLen = 1
	}

	var prefixBuf bytes.Buffer
	prefixBytes := lineText[:pos.ByteColumn-1]

	// Iterate character-by-character to properly pad spacing,
	// preserving the physical layout of tabs vs. spaces in the user's terminal.
	for len(prefixBytes) > 0 {
		r, size := utf8.DecodeRune(prefixBytes)
		if r == '\t' {
			prefixBuf.WriteByte('\t')
		} else {
			prefixBuf.WriteByte(' ')
		}
		prefixBytes = prefixBytes[size:]
	}

	squiggles := strings.Repeat("^", squigglyLen)
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, "%s\n", loc.String())
	_, _ = fmt.Fprintf(&buf, "%4d | %s\n", pos.Line, lineText)
	_, _ = fmt.Fprintf(&buf, "     | %s%s\n", prefixBuf.String(), squiggles)

	return buf.String()
}
