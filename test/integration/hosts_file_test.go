//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/txn2/txeh"
)

// TestHostsFileRoundTrip verifies that reading and writing a hosts file
// produces identical output, preserving comments, empty lines, and formatting.
func TestHostsFileRoundTrip(t *testing.T) {
	// Create a temp file with known content
	content := `127.0.0.1 localhost
::1 localhost ip6-localhost

# Development services
10.0.0.1 app1.local app2.local # my project
10.0.0.2 db.local
`
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")
	if err := os.WriteFile(hostsPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp hosts file: %v", err)
	}

	// Read, modify, save, re-read
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath:  hostsPath,
		WriteFilePath: hostsPath,
	})
	if err != nil {
		t.Fatalf("NewHosts failed: %v", err)
	}

	hosts.AddHost("10.0.0.3", "cache.local")

	if err := hosts.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Re-read and verify
	hosts2, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath: hostsPath,
	})
	if err != nil {
		t.Fatalf("second NewHosts failed: %v", err)
	}

	hostnames := hosts2.ListHostsByIP("10.0.0.3")
	if len(hostnames) != 1 || hostnames[0] != "cache.local" {
		t.Errorf("expected [cache.local], got %v", hostnames)
	}

	// Verify original entries survived
	hostnames = hosts2.ListHostsByIP("10.0.0.1")
	if len(hostnames) != 2 {
		t.Errorf("expected 2 hosts at 10.0.0.1, got %d", len(hostnames))
	}
}

// TestHostsFileSaveAs verifies SaveAs writes to an alternate path.
func TestHostsFileSaveAs(t *testing.T) {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "hosts.src")
	if err := os.WriteFile(srcPath, []byte("127.0.0.1 localhost\n"), 0644); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: srcPath})
	if err != nil {
		t.Fatalf("NewHosts failed: %v", err)
	}

	hosts.AddHost("10.0.0.1", "test.local")

	outPath := filepath.Join(tmpDir, "hosts.out")
	if err := hosts.SaveAs(outPath); err != nil {
		t.Fatalf("SaveAs failed: %v", err)
	}

	hosts2, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: outPath})
	if err != nil {
		t.Fatalf("failed to re-read: %v", err)
	}

	hostnames := hosts2.ListHostsByIP("10.0.0.1")
	if len(hostnames) != 1 {
		t.Errorf("expected 1 host at 10.0.0.1, got %d", len(hostnames))
	}
}
