package codegen

import "fmt"

func (t *Transpiler) write(s string) {
	t.buf.WriteString(s)
}

func (t *Transpiler) printf(format string, args ...any) {
	fmt.Fprintf(&t.buf, format, args...)
}

func (t *Transpiler) newline() {
	t.buf.WriteByte('\n')
}
