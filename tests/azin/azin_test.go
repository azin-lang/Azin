package azin_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var compiler = os.Getenv("AZC")

func TestAzinPrograms(t *testing.T) {
	compiler := os.Getenv("AZC")

	if runtime.GOOS == "windows" &&
		compiler != "" &&
		filepath.Ext(compiler) == "" {
		compiler += ".exe"
	}

	runDir(t, "success", false)
	runDir(t, "fail", true)
}

func runDir(t *testing.T, dir string, shouldFail bool) {
	files, err := filepath.Glob(filepath.Join(dir, "*.az"))
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		file := file

		t.Run(filepath.Base(file), func(t *testing.T) {
			defer func() {
				_ = os.Remove("output.exe")
			}()
			cmd := exec.Command(compiler, file)

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()

			if shouldFail {
				if err == nil {
					t.Fatalf("expected compilation to fail")
				}

				expected := expectedError(file)
				if expected != "" &&
					!strings.Contains(out.String(), expected) {

					t.Fatalf(
						"expected error containing:\n\n%s\n\nrun error: %v\n\noutput:\n%s",
						expected,
						err,
						out.String(),
					)
				}
			} else {
				if err != nil {
					t.Fatalf(
						"expected compilation success\n\nrun error: %v\n\noutput:\n%s",
						err,
						out.String(),
					)
				}
			}
		})
	}
}

func expectedError(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		return ""
	}

	var b strings.Builder

	lines := strings.Split(string(data), "\n")

	found := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if !found {
			if line == "// error expected:" {
				found = true
			}
			continue
		}

		if !strings.HasPrefix(line, "//") {
			break
		}

		b.WriteString(strings.TrimSpace(strings.TrimPrefix(line, "//")))
		b.WriteByte('\n')
	}

	return strings.TrimSpace(b.String())
}
