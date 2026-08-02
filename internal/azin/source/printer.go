package source

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// PrintColoredDiagnostic generates a fully color-coordinated CLI diagnostic block.
func PrintColoredDiagnostic(loc Location, label, note, help string, pipeColor, styleColor func(a ...any) string) string {
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

	for len(prefixBytes) > 0 {
		r, size := utf8.DecodeRune(prefixBytes)
		if r == '\t' {
			prefixBuf.WriteByte('\t')
		} else {
			prefixBuf.WriteByte(' ')
		}
		prefixBytes = prefixBytes[size:]
	}

	var squiggles string
	if squigglyLen <= 1 {
		squiggles = "^"
	} else {
		squiggles = "^" + strings.Repeat("~", squigglyLen-1)
	}

	var buf bytes.Buffer

	// Line number and source line
	buf.WriteString(fmt.Sprintf("%s | %s\n", pipeColor(fmt.Sprintf("%4d", pos.Line)), lineText))

	// Underline and label
	if label != "" {
		buf.WriteString(fmt.Sprintf("%s | %s%s %s\n", pipeColor("    "), prefixBuf.String(), styleColor(squiggles), styleColor(label)))
	} else {
		buf.WriteString(fmt.Sprintf("%s | %s%s\n", pipeColor("    "), prefixBuf.String(), styleColor(squiggles)))
	}

	// Empty separator pipe if we have footer annotations
	if note != "" || help != "" {
		buf.WriteString(fmt.Sprintf("%s |\n", pipeColor("    ")))
	}

	// Actionable note or suggestion
	if note != "" {
		buf.WriteString(fmt.Sprintf("%s note: %s\n", pipeColor("    ="), note))
	}

	// Architectural help guidance
	if help != "" {
		buf.WriteString(fmt.Sprintf("%s help: %s\n", pipeColor("    ="), help))
	}

	return buf.String()
}
