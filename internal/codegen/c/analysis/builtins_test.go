package analysis

import (
	"testing"
)

func TestLookupBuiltin(t *testing.T) {
	tests := []struct {
		name          string
		exceptedValue string
		exceptedOk    bool
	}{
		{name: "abs", exceptedValue: Stdlib, exceptedOk: true},
		{name: "not_function", exceptedValue: "", exceptedOk: false},
		{name: "printf", exceptedValue: Stdio, exceptedOk: true},
		{name: "atoi", exceptedValue: Stdlib, exceptedOk: true},
		{name: "", exceptedValue: "", exceptedOk: false},
		{name: "ATOF", exceptedValue: "", exceptedOk: false},
		{name: " malloc ", exceptedValue: "", exceptedOk: false},
	}

	for _, tt := range tests {
		result, ok := LookupBuiltin(tt.name)
		if result.Include != tt.exceptedValue || ok != tt.exceptedOk {
			t.Errorf("LookupBuiltin(%q) = (%q, %v); want (%q, %v)",
				tt.name, result, ok, tt.exceptedValue, tt.exceptedOk)
		}
	}
}
