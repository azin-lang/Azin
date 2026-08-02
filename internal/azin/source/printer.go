package source

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

type diagnosticLayout struct {
	Location string

	LineNumber uint32
	LineText   []byte

	PointerIndent string
	PointerWidth  int
}

// PrintColoredDiagnostic renders a colored diagnostic.
func PrintColoredDiagnostic(
	loc Location,
	label, note, help string,
	width int,
	pipe,
	style,
	bold func(...any) string,
) string {
	if loc.File == nil || !loc.Span.IsValidAndNonEmpty() {
		return "<invalid location>\n"
	}

	layout := makeDiagnosticLayout(loc)

	var buf bytes.Buffer

	writeDiagnosticHeader(&buf, layout, pipe)
	writeDiagnosticSource(&buf, layout, label, pipe, style)
	writeDiagnosticFooter(&buf, note, help, width, bold)

	return buf.String()
}

func makeDiagnosticLayout(loc Location) diagnosticLayout {
	pos := loc.Position()
	line := loc.File.Line(pos.Line)
	text := loc.File.BytesOf(NewSpan(line.Start, line.End))

	column := clampColumn(pos.ByteColumn, len(text))

	endColumn := pos.ByteColumn + uint32(max(pointerWidth(loc, text, column)-1, 0))

	return diagnosticLayout{
		Location:      formatLocation(loc, pos, endColumn),
		LineNumber:    pos.Line,
		LineText:      text,
		PointerIndent: pointerIndent(text[:column]),
		PointerWidth:  pointerWidth(loc, text, column),
	}
}

func writeDiagnosticHeader(
	buf *bytes.Buffer,
	layout diagnosticLayout,
	pipe func(...any) string,
) {
	fmt.Fprintf(buf, "   %s %s\n", pipe("┌─"), pipe(layout.Location))
	fmt.Fprintf(buf, "   %s\n", pipe("│"))
}

func writeDiagnosticSource(
	buf *bytes.Buffer,
	layout diagnosticLayout,
	label string,
	pipe func(...any) string,
	style func(...any) string,
) {
	pointer := strings.Repeat("╍", layout.PointerWidth)

	fmt.Fprintf(buf, "%2d %s %s\n",
		layout.LineNumber,
		pipe("│"),
		layout.LineText,
	)

	fmt.Fprintf(buf, "   %s %s%s\n",
		pipe("│"),
		layout.PointerIndent,
		style(pointer),
	)

	if label != "" {
		fmt.Fprintf(buf, "   %s %s%s\n",
			pipe("│"),
			layout.PointerIndent,
			style(label),
		)
	}

	fmt.Fprintf(buf, "   %s\n\n", pipe("└"))
}

func writeDiagnosticFooter(
	buf *bytes.Buffer,
	note,
	help string,
	width int,
	bold func(...any) string,
) {
	if note != "" {
		buf.WriteString(wrapText(note, width))
		buf.WriteByte('\n')
	}

	if help == "" {
		return
	}

	if note != "" {
		buf.WriteByte('\n')
	}

	for _, line := range strings.Split(wrapText(help, width), "\n") {
		buf.WriteString(bold(line))
		buf.WriteByte('\n')
	}
}

func formatLocation(
	loc Location,
	pos Position,
	endColumn uint32,
) string {
	if pos.ByteColumn == endColumn {
		return fmt.Sprintf(
			"%s:%d:%d",
			loc.File.Path(),
			pos.Line,
			pos.ByteColumn,
		)
	}

	return fmt.Sprintf(
		"%s:%d:%d-%d",
		loc.File.Path(),
		pos.Line,
		pos.ByteColumn,
		endColumn,
	)
}

func clampColumn(column uint32, lineLength int) int {
	column--

	if int(column) > lineLength {
		return lineLength
	}

	return max(int(column), 0)
}

func pointerWidth(loc Location, line []byte, column int) int {
	width := utf8.RuneCount(loc.File.BytesOf(loc.Span))

	if width <= 1 {
		width = scanPointerWidth(line[column:])
	}

	return min(max(width, 1), max(len(line)-column, 1))
}

func scanPointerWidth(line []byte) int {
	width := 0

	for len(line) > 0 {
		r, size := utf8.DecodeRune(line)

		switch r {
		case ' ', '\t', '\n', '\r', ';', ',', '(', ')', '[', ']', '{', '}':
			return max(width, 1)
		}

		width++
		line = line[size:]
	}

	return max(width, 1)
}

func pointerIndent(line []byte) string {
	var b strings.Builder

	for len(line) > 0 {
		r, size := utf8.DecodeRune(line)

		if r == '\t' {
			b.WriteByte('\t')
		} else {
			b.WriteByte(' ')
		}

		line = line[size:]
	}

	return b.String()
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var b strings.Builder

	for i, paragraph := range strings.Split(text, "\n") {
		if i > 0 {
			b.WriteByte('\n')
		}

		column := 0

		for _, word := range strings.Fields(paragraph) {
			length := utf8.RuneCountInString(word)

			if column != 0 && column+1+length > width {
				b.WriteByte('\n')
				column = 0
			}

			if column != 0 {
				b.WriteByte(' ')
				column++
			}

			b.WriteString(word)
			column += length
		}
	}

	return b.String()
}
