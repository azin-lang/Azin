// Package fs implements file operations for Azin source files.
package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SourceExtension is the expected file extension for source files.
const SourceExtension = ".az"

func ValidateSourcePath(path string) (string, error) {
	path = filepath.Clean(path)

	base := filepath.Base(path)
	if strings.EqualFold(base, SourceExtension) {
		return "", fmt.Errorf("invalid source file %q: missing file name", path)
	}

	if !strings.EqualFold(filepath.Ext(base), SourceExtension) {
		return "", fmt.Errorf(
			"invalid source file %q: expected %q extension",
			path,
			SourceExtension,
		)
	}

	return path, nil
}

// ReadSourceFile reads the contents of a source file.
func ReadSourceFile(path string, ignoreExtension bool) ([]byte, error) {
	var err error

	path = filepath.Clean(path)

	if !ignoreExtension {
		path, err = ValidateSourcePath(path)
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			return nil, fmt.Errorf("source file %q does not exist", path)
		case errors.Is(err, os.ErrPermission):
			return nil, fmt.Errorf("permission denied reading %q", path)
		case errors.Is(err, os.ErrInvalid):
			return nil, fmt.Errorf("invalid source file %q", path)
		default:
			if info, statErr := os.Stat(path); statErr == nil && info.IsDir() {
				return nil, fmt.Errorf("%q is a directory", path)
			}

			return nil, fmt.Errorf("read %q: %w", path, err)
		}
	}

	return data, nil
}
