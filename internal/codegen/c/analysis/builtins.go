package analysis

type Builtin struct {
	Include string
}

var Builtins = map[string]Builtin{
	"printf":   {"stdio.h"},
	"fprintf":  {"stdio.h"},
	"sprintf":  {"stdio.h"},
	"snprintf": {"stdio.h"},
	"scanf":    {"stdio.h"},
	"sscanf":   {"stdio.h"},

	"malloc":  {"stdlib.h"},
	"calloc":  {"stdlib.h"},
	"realloc": {"stdlib.h"},
	"free":    {"stdlib.h"},
	"exit":    {"stdlib.h"},
	"abs":     {"stdlib.h"},

	"strlen": {"string.h"},
	"strcpy": {"string.h"},
	"strcmp": {"string.h"},
	"memset": {"string.h"},
	"memcpy": {"string.h"},
}

func LookupBuiltin(name string) (Builtin, bool) {
	v, ok := Builtins[name]
	return v, ok
}
