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
	"fopen":    {"stdio.h"},
	"fclose":   {"stdio.h"},
	"fread":    {"stdio.h"},
	"fwrite":   {"stdio.h"},
	"fflush":   {"stdio.h"},
	"fseek":    {"stdio.h"},
	"ftell":    {"stdio.h"},
	"puts":     {"stdio.h"},
	"fgets":    {"stdio.h"},
	"fputs":    {"stdio.h"},
	"getchar":  {"stdio.h"},
	"putchar":  {"stdio.h"},

	"malloc":  {"stdlib.h"},
	"calloc":  {"stdlib.h"},
	"realloc": {"stdlib.h"},
	"free":    {"stdlib.h"},
	"exit":    {"stdlib.h"},
	"abs" :    {"stdlib.h"},
	"atoi" :   {"stdlib.h"},
	"atof" :   {"stdlib.h"},
	"atol" :   {"stdlib.h"},
	"atoll" :  {"stdlib.h"},
	"strtol":  {"stdlib.h"},
	"strtoll": {"stdlib.h"},
	"strtoul": {"stdlib.h"},
	"strtoull":{"stdlib.h"},
	"strtod":  {"stdlib.h"},
	"strtof":  {"stdlib.h"},

	"strlen":  {"string.h"},
	"strcpy":  {"string.h"},
	"strcmp":  {"string.h"},
	"memset":  {"string.h"},
	"memcpy":  {"string.h"},
	"strncpy": {"string.h"},
	"strcat":  {"string.h"},
	"strncat": {"string.h"},
	"strncmp": {"string.h"},
	"strchr":  {"string.h"},
	"strrchr": {"string.h"},
	"strstr":  {"string.h"},
	"strtok":  {"string.h"},
	"strdup":  {"string.h"},
	"memmove": {"string.h"},
	"memcmp":  {"string.h"},
}

func LookupBuiltin(name string) (Builtin, bool) {
	v, ok := Builtins[name]
	return v, ok
}
