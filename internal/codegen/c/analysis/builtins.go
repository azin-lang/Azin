package analysis

type Builtin struct {
	Include string
}

const (
	Stdio  = "stdio.h"
	Stdlib = "stdlib.h"
	String = "string.h"
)

var Builtins = map[string]Builtin{
	"printf":   {Stdio},
	"fprintf":  {Stdio},
	"sprintf":  {Stdio},
	"snprintf": {Stdio},
	"scanf":    {Stdio},
	"sscanf":   {Stdio},
	"fopen":    {Stdio},
	"fclose":   {Stdio},
	"fread":    {Stdio},
	"fwrite":   {Stdio},
	"fflush":   {Stdio},
	"fseek":    {Stdio},
	"ftell":    {Stdio},
	"puts":     {Stdio},
	"fgets":    {Stdio},
	"fputs":    {Stdio},
	"getchar":  {Stdio},
	"putchar":  {Stdio},

	"malloc":  {Stdlib},
	"calloc":  {Stdlib},
	"realloc": {Stdlib},
	"free":    {Stdlib},
	"exit":    {Stdlib},
	"abs" :    {Stdlib},
	"atoi" :   {Stdlib},
	"atof" :   {Stdlib},
	"atol" :   {Stdlib},
	"atoll" :  {Stdlib},
	"strtol":  {Stdlib},
	"strtoll": {Stdlib},
	"strtoul": {Stdlib},
	"strtoull":{Stdlib},
	"strtod":  {Stdlib},
	"strtof":  {Stdlib},

	"strlen":  {String},
	"strcpy":  {String},
	"strcmp":  {String},
	"memset":  {String},
	"memcpy":  {String},
	"strncpy": {String},
	"strcat":  {String},
	"strncat": {String},
	"strncmp": {String},
	"strchr":  {String},
	"strrchr": {String},
	"strstr":  {String},
	"strtok":  {String},
	"memmove": {String},
	"memcmp":  {String},
}

func LookupBuiltin(name string) (Builtin, bool) {
	v, ok := Builtins[name]
	return v, ok
}
