package c

import (
	"bytes"
	"fmt"
)

type writer struct {
	buf    bytes.Buffer
	indent int
	err    error
}

func (w *writer) reset() {
	w.buf.Reset()
	w.indent = 0
	w.err = nil
}

func (w *writer) String() string {
	return w.buf.String()
}

func (w *writer) write(s string) {
	if w.err != nil {
		return
	}
	_, w.err = w.buf.WriteString(s)
}

func (w *writer) printf(
	format string,
	args ...any,
) {
	if w.err != nil {
		return
	}
	_, w.err = fmt.Fprintf(&w.buf, format, args...)
}

func (w *writer) newline() {
	if w.err != nil {
		return
	}
	w.err = w.buf.WriteByte('\n')
}

func (w *writer) indentLine() {
	for range w.indent {
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
