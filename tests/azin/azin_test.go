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

var compiler = func() string {
	c := os.Getenv("AZC")

	if runtime.GOOS == "windows" &&
		c != "" &&
		filepath.Ext(c) == "" {
		c += ".exe"
	}

	return c
}()

func TestAzinPrograms(t *testing.T) {
	if compiler == "" {
		t.Fatal("AZC environment variable is not set")
	}

	runDir(t, "success", false)
	runDir(t, "fail", true)
}

func runDir(t *testing.T, dir string, shouldFail bool) {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(dir, "*.az"))
	if err != nil {
		t.Fatalf("finding test files: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("no test files found in %q", dir)
	}

	for _, file := range files {

		t.Run(filepath.Base(file), func(t *testing.T) {
			defer os.Remove("output.exe")

			cmd := exec.Command(compiler, file)

			var output bytes.Buffer
			cmd.Stdout = &output
			cmd.Stderr = &output

			err := cmd.Run()

			got := strings.TrimSpace(output.String())

			if shouldFail {
				if err == nil {
					t.Fatal("expected compilation to fail")
				}

				expected := expectedError(t, file)
				if expected == "" {
					t.Fatal("missing '// error expected:' block")
				}

				if !strings.Contains(got, expected) {
					t.Fatalf(
						"unexpected compiler output\n\nexpected to contain:\n%s\n\ngot:\n%s",
						expected,
						got,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"compilation failed:\n\n%s",
					got,
				)
			}
		})
	}
}

func expectedError(t *testing.T, file string) string {
	t.Helper()

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("reading %q: %v", file, err)
	}

	lines := strings.Split(string(data), "\n")

	var b strings.Builder
	found := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if !found {
			found = line == "// error expected:"
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
