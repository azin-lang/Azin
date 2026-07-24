package analysis

type Builtin struct {
	Include string
}

const (
	StdioHeader  = "stdio.h"
	StdlibHeader = "stdlib.h"
	StringHeader = "string.h"
)

var Builtins = map[string]Builtin{
	"printf":   {StdioHeader},
	"fprintf":  {StdioHeader},
	"sprintf":  {StdioHeader},
	"snprintf": {StdioHeader},
	"scanf":    {StdioHeader},
	"sscanf":   {StdioHeader},

	"malloc":  {StdlibHeader},
	"calloc":  {StdlibHeader},
	"realloc": {StdlibHeader},
	"free":    {StdlibHeader},
	"exit":    {StdlibHeader},
	"abs":     {StdlibHeader},

	"strlen": {StringHeader},
	"strcpy": {StringHeader},
	"strcmp": {StringHeader},
	"memset": {StringHeader},
	"memcpy": {StringHeader},
}

func LookupBuiltin(name string) (Builtin, bool) {
	v, ok := Builtins[name]
	return v, ok
}
