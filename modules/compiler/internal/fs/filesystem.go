// Package fs implements file operations for Azin source files.
package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// SourceExtension is the expected file extension for source files.
const SourceExtension = ".az"

func validateSourceFile(path string) error {
	if ext := filepath.Ext(path); ext != SourceExtension {
		return fmt.Errorf("invalid source file %q: expected %q extension, got %q", path, SourceExtension, ext)
	}
	return nil
}

// ReadSourceFile reads the contents of the given file path.
func ReadSourceFile(path string, ignoreExtension bool) ([]byte, error) {
	if !ignoreExtension {
		if err := validateSourceFile(path); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}

	return data, nil
}
