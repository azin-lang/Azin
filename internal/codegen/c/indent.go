package c

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
