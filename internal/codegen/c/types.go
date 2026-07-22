package c

func emitType(name string) string {
	switch name {
	case "unit":
		return "void"

	case "int":
		return "int"

	case "float":
		return "float"

	case "char":
		return "char"

	case "string":
		return "const char *"

	case "bool":
		return "bool"

	default:
		// assume it's a user defined type
		// TODO: check if it is, instead of assuming and letting C handle it
		return name
	}
}
