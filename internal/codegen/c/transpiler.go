package c

import "bytes"

type Transpiler struct {
	buf bytes.Buffer

	indent int

	enums map[string]bool
}
