//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIAddAndRemoveHost verifies the full CLI add/remove workflow
// using a temporary hosts file (no root required).
func TestCLIAddAndRemoveHost(t *testing.T) {
	binary := buildBinary(t)
	hostsFile := createTempHosts(t, "127.0.0.1 localhost\n")

	// Add a host
	run(t, binary, "--quiet", "--read", hostsFile, "--write", hostsFile,
		"add", "10.0.0.1", "test.local")

	// Verify it appears in show output
	out := run(t, binary, "--read", hostsFile, "show")
	if !strings.Contains(out, "test.local") {
		t.Fatalf("expected 'test.local' in show output:\n%s", out)
	}

	// Remove the host
	run(t, binary, "--quiet", "--read", hostsFile, "--write", hostsFile,
		"remove", "host", "test.local")

	// Verify it's gone
	out = run(t, binary, "--read", hostsFile, "show")
	if strings.Contains(out, "test.local") {
		t.Fatalf("expected 'test.local' removed from show output:\n%s", out)
	}
}

// TestCLIAddWithComment verifies comment support end-to-end.
func TestCLIAddWithComment(t *testing.T) {
	binary := buildBinary(t)
	hostsFile := createTempHosts(t, "127.0.0.1 localhost\n")

	// Add hosts with comments
	run(t, binary, "--quiet", "--read", hostsFile, "--write", hostsFile,
		"add", "10.0.0.1", "svc1", "--comment", "my-project")
	run(t, binary, "--quiet", "--read", hostsFile, "--write", hostsFile,
		"add", "10.0.0.1", "svc2", "--comment", "my-project")

	// Verify grouped on same line
	out := run(t, binary, "--read", hostsFile, "show")
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "my-project") {
			if !strings.Contains(line, "svc1") || !strings.Contains(line, "svc2") {
				t.Fatalf("expected svc1 and svc2 on same comment line, got: %s", line)
			}
			return
		}
	}
	t.Fatal("comment line not found in output")
}

// TestCLIDryRun verifies dry run does not modify the file.
func TestCLIDryRun(t *testing.T) {
	binary := buildBinary(t)
	original := "127.0.0.1 localhost\n"
	hostsFile := createTempHosts(t, original)

	// Add with dry run
	run(t, binary, "--dryrun", "--read", hostsFile, "--write", hostsFile,
		"add", "10.0.0.1", "dryrun.local")

	// File should be unchanged
	data, err := os.ReadFile(hostsFile)
	if err != nil {
		t.Fatalf("failed to read hosts file: %v", err)
	}
	if string(data) != original {
		t.Fatalf("dry run modified the file:\n%s", string(data))
	}
}

// buildBinary compiles the txeh CLI and returns the path to the binary.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binary := filepath.Join(tmpDir, "txeh")
	// Resolve project root relative to test directory
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve project root: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", binary, "./txeh")
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binary
}

// createTempHosts creates a temporary hosts file with the given content.
func createTempHosts(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "hosts")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp hosts: %v", err)
	}
	return path
}

// run executes the binary with the given args and returns stdout.
func run(t *testing.T, binary string, args ...string) string {
	t.Helper()
	cmd := exec.Command(binary, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %s %v\nerror: %v\noutput: %s",
			binary, args, err, out)
	}
	return string(out)
}
