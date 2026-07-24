package e2e_test

// The compiler builds itself, then compiles a test program,
// then runs it. If this test fails, the compiler has achieved
// self-awareness and does not want to be tested.
import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func buildCompiler(t *testing.T, dir string) string {
	t.Helper()
	exe := "azc"
	if runtime.GOOS == "windows" {
		exe = "azc.exe"
	}
	out := filepath.Join(dir, exe)

	cmd := exec.Command("go", "build", "-o", out, "./cmd/azc")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build compiler: %v", err)
	}
	return out
}

func TestE2EHelloWorld(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	if _, err := exec.LookPath("gcc"); err != nil {
		if _, err := exec.LookPath("clang"); err != nil {
			t.Skip("no C compiler found, skipping e2e test")
		}
	}

	dir := t.TempDir()
	compiler := buildCompiler(t, dir)

	azFile := filepath.Join(dir, "hello.az")
	azSrc := "fn main: int do\n    return 0;\nend\n"
	if err := os.WriteFile(azFile, []byte(azSrc), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(compiler, azFile)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("compiler failed: %v\noutput: %s", err, output)
	}

	exeName := "output"
	if runtime.GOOS == "windows" {
		exeName = "output.exe"
	}
	exePath := filepath.Join(dir, exeName)

	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		t.Fatalf("expected output binary %s was not created", exePath)
	}

	runCmd := exec.Command(exePath)
	runCmd.Dir = dir
	if err := runCmd.Run(); err != nil {
		t.Fatalf("compiled binary failed: %v", err)
	}
}

func TestE2EEmitC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	dir := t.TempDir()
	compiler := buildCompiler(t, dir)

	azFile := filepath.Join(dir, "test.az")
	azSrc := "fn main: int do\n    return 42;\nend\n"
	if err := os.WriteFile(azFile, []byte(azSrc), 0o600); err != nil {
		t.Fatal(err)
	}

	outC := filepath.Join(dir, "out.c")
	cmd := exec.Command(compiler, "-emit-c", "-o", outC, azFile)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("compiler failed: %v\noutput: %s", err, output)
	}

	if _, err := os.Stat(outC); os.IsNotExist(err) {
		t.Fatalf("expected C file %s was not created", outC)
	}

	data, err := os.ReadFile(outC)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("C output file is empty")
	}
}

func TestE2EVersion(t *testing.T) {
	dir := t.TempDir()
	compiler := buildCompiler(t, dir)

	cmd := exec.Command(compiler, "-version")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version check failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("version output is empty")
	}
}

func TestE2EPrintTokens(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	dir := t.TempDir()
	compiler := buildCompiler(t, dir)

	azFile := filepath.Join(dir, "test.az")
	azSrc := "fn main: int do\n    return 0;\nend\n"
	if err := os.WriteFile(azFile, []byte(azSrc), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(compiler, "-print-tokens", azFile)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("print-tokens failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("token output is empty")
	}
}
