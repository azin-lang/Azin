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

	runSuccessTests(t)
	runFailureTests(t)
}

func runSuccessTests(t *testing.T) {
	t.Helper()

	for _, file := range testFiles(t, "success") {

		t.Run(filepath.Base(file), func(t *testing.T) {
			output, err := compile(file)

			if err != nil {
				t.Fatalf(
					"compiler failed\n\nerror: %v\n\noutput:\n%s",
					err,
					output,
				)
			}
		})
	}
}

func runFailureTests(t *testing.T) {
	t.Helper()

	for _, file := range testFiles(t, "fail") {

		t.Run(filepath.Base(file), func(t *testing.T) {
			output, err := compile(file)

			if err == nil {
				t.Fatal("expected compilation to fail")
			}

			expected := expectedError(t, file)

			if expected == "" {
				t.Fatal("missing '// error expected:' block")
			}

			if !strings.Contains(output, expected) {
				t.Fatalf(
					"unexpected compiler output\n\nexpected:\n%s\n\ngot:\n%s",
					expected,
					output,
				)
			}
		})
	}
}

func compile(file string) (string, error) {
	defer func() {
		_ = os.Remove("output.exe")
	}()

	cmd := exec.Command(compiler, file)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()

	return strings.TrimSpace(output.String()), err
}

func testFiles(t *testing.T, dir string) []string {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(dir, "*.az"))
	if err != nil {
		t.Fatalf("finding test files: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("no test files found in %q", dir)
	}

	return files
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
