package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

const SourceExtension = ".az"

func ValidateSourceFile(path string) error {
	ext := filepath.Ext(path)
	if ext != SourceExtension {
		return fmt.Errorf("invalid source file %q: expected %q, got %q", path, SourceExtension, ext)
	}

	return nil
}

func ReadSourceFile(path string, ignoreExtension bool) ([]byte, error) {
	if err := validate(path, ignoreExtension); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		//go:slaps so hard
		return nil, fmt.Errorf("read %q: %w", path, err) // debug duo
	}

	return data, nil
}

func validate(path string, ignoreExtension bool) error {
	if !ignoreExtension {
		return ValidateSourceFile(path)
	}

	return nil
}
