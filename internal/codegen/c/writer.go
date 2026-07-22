package c

import "fmt"

func (t *Transpiler) write(s string) {
	t.buf.WriteString(s)
}

func (t *Transpiler) printf(format string, args ...any) {
	_, err := fmt.Fprintf(&t.buf, format, args...)
	if err != nil {
		return
	}
}

func (t *Transpiler) newline() {
	t.buf.WriteByte('\n')
}

func (t *Transpiler) writeIndent() {
	for i := 0; i < t.indent; i++ {
		t.write("    ")
	}
}

func (t *Transpiler) pushIndent() {
	t.indent++
}

func (t *Transpiler) popIndent() {
	t.indent--
}
