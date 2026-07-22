package c

import (
	"bytes"
	"fmt"
)

type writer struct {
	buf    bytes.Buffer
	indent int
}

func (w *writer) reset() {
	w.buf.Reset()
	w.indent = 0
}

func (w *writer) String() string {
	return w.buf.String()
}

func (w *writer) write(s string) {
	_, _ = w.buf.WriteString(s)
}

func (w *writer) printf(
	format string,
	args ...any,
) {
	_, _ = fmt.Fprintf(&w.buf, format, args...)
}

func (w *writer) newline() {
	_ = w.buf.WriteByte('\n')
}

func (w *writer) indentLine() {
	for i := 0; i < w.indent; i++ {
		w.write("\t")
	}
}

func (w *writer) pushIndent() {
	w.indent++
}

func (w *writer) popIndent() {
	if w.indent > 0 {
		w.indent--
	}
}
