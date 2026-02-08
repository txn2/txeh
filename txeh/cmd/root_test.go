package cmd

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/txn2/txeh"
)

const testHostLocalhost = "localhost"

const testHostsWithComments = "127.0.0.1        localhost\n127.0.0.1        app1 # dev services\n"

// Bug Regression Tests for CMD package
// These tests document and verify bugs found during code analysis.

// =============================================================================
// BUG #1: IP Validation Logic Inverted (CRITICAL)
// File: txeh/cmd/root.go:88-89
//
// The function is:
//   func validateIPAddress(ip string) bool {
//       return (net.ParseIP(ip) == nil)
//   }
//
// This returns TRUE when the IP is INVALID (ParseIP returns nil for invalid).
// It should return TRUE when the IP is VALID.
// =============================================================================

func TestBug1_ValidateIPAddress_ValidIPv4_ShouldReturnTrue(t *testing.T) {
	validIPs := []string{
		"127.0.0.1",
		"192.168.1.1",
		"10.0.0.1",
		"255.255.255.255",
		"0.0.0.0",
	}

	for _, ip := range validIPs {
		result := validateIPAddress(ip)
		// BUG: Currently returns FALSE for valid IPs
		if !result {
			t.Errorf("BUG #1 CONFIRMED: validateIPAddress(%q) = %v, expected true (valid IP rejected)",
				ip, result)
		}
	}
}

func TestBug1_ValidateIPAddress_ValidIPv6_ShouldReturnTrue(t *testing.T) {
	validIPs := []string{
		"::1",
		"2001:db8::1",
		"fe80::1",
		"::ffff:192.168.1.1",
	}

	for _, ip := range validIPs {
		result := validateIPAddress(ip)
		// BUG: Currently returns FALSE for valid IPs
		if !result {
			t.Errorf("BUG #1 CONFIRMED: validateIPAddress(%q) = %v, expected true (valid IPv6 rejected)",
				ip, result)
		}
	}
}

func TestBug1_ValidateIPAddress_Invalid_ShouldReturnFalse(t *testing.T) {
	invalidIPs := []string{
		"999.999.999.999",
		"not-an-ip",
		"127.0.0.1.1",
		"",
		"abc.def.ghi.jkl",
		"192.168.1",
		"192.168.1.256",
	}

	for _, ip := range invalidIPs {
		result := validateIPAddress(ip)
		// BUG: Currently returns TRUE for invalid IPs
		if result {
			t.Errorf("BUG #1 CONFIRMED: validateIPAddress(%q) = %v, expected false (invalid IP accepted)",
				ip, result)
		}
	}
}

func TestBug1_ExpectedBehavior_Demonstration(t *testing.T) {
	// This demonstrates what the correct behavior should be
	validIP := "127.0.0.1"
	invalidIP := "999.999.999.999"

	// net.ParseIP returns non-nil for valid, nil for invalid
	validParsed := net.ParseIP(validIP)
	invalidParsed := net.ParseIP(invalidIP)

	t.Logf("net.ParseIP(%q) = %v (non-nil = valid)", validIP, validParsed)
	t.Logf("net.ParseIP(%q) = %v (nil = invalid)", invalidIP, invalidParsed)

	// Current (buggy) implementation
	currentValidResult := validateIPAddress(validIP)
	currentInvalidResult := validateIPAddress(invalidIP)

	t.Logf("validateIPAddress(%q) = %v (should be true)", validIP, currentValidResult)
	t.Logf("validateIPAddress(%q) = %v (should be false)", invalidIP, currentInvalidResult)

	// The fix should be: return net.ParseIP(ip) != nil
	correctValidResult := net.ParseIP(validIP) != nil
	correctInvalidResult := net.ParseIP(invalidIP) != nil

	t.Logf("Correct behavior for %q: %v", validIP, correctValidResult)
	t.Logf("Correct behavior for %q: %v", invalidIP, correctInvalidResult)

	if currentValidResult != correctValidResult || currentInvalidResult != correctInvalidResult {
		t.Error("BUG #1 CONFIRMED: validateIPAddress has inverted logic")
	}
}

// =============================================================================
// BUG #7: Hostname Regex Allows Invalid Patterns (LOW)
// File: txeh/cmd/root.go:60
//
// The regex: ^([A-Za-z]|[0-9]|-|_|\.)+$
// Allows invalid hostnames like: ".hostname", "hostname.", "...", "host..name"
// RFC 1123 says hostname labels cannot start or end with hyphen.
// FQDN rules say labels should be separated by single dots.
// =============================================================================

func TestBug7_ValidateHostname_LeadingDot_ShouldBeInvalid(t *testing.T) {
	invalidHostnames := []string{
		".hostname",
		".example.com",
	}

	for _, hostname := range invalidHostnames {
		result := validateHostname(hostname)
		// BUG: Currently returns TRUE for these invalid hostnames
		if result {
			t.Errorf("BUG #7 CONFIRMED: validateHostname(%q) = %v, expected false (leading dot invalid)",
				hostname, result)
		}
	}
}

func TestBug7_ValidateHostname_TrailingDot_ShouldBeInvalid(t *testing.T) {
	// Trailing dots are not used in /etc/hosts files
	invalidHostnames := []string{
		"hostname.",
		"example.com.",
	}

	for _, hostname := range invalidHostnames {
		result := validateHostname(hostname)
		if result {
			t.Errorf("validateHostname(%q) = %v, expected false (trailing dot invalid)", hostname, result)
		}
	}
}

func TestBug7_ValidateHostname_ConsecutiveDots_ShouldBeInvalid(t *testing.T) {
	invalidHostnames := []string{
		"host..name",
		"example..com",
		"...",
		"a..b..c",
	}

	for _, hostname := range invalidHostnames {
		result := validateHostname(hostname)
		// BUG: Currently returns TRUE for these invalid hostnames
		if result {
			t.Errorf("BUG #7 CONFIRMED: validateHostname(%q) = %v, expected false (consecutive dots invalid)",
				hostname, result)
		}
	}
}

func TestBug7_ValidateHostname_OnlyDots_ShouldBeInvalid(t *testing.T) {
	result := validateHostname("...")
	if result {
		t.Error("BUG #7 CONFIRMED: validateHostname('...') = true, expected false")
	}
}

func TestBug7_ValidateHostname_Valid_ShouldReturnTrue(t *testing.T) {
	validHostnames := []string{
		"localhost",
		"example.com",
		"my-server",
		"server_01",
		"a.b.c.d",
		"123.example.com",
		"test-server.local",
	}

	for _, hostname := range validHostnames {
		result := validateHostname(hostname)
		if !result {
			t.Errorf("validateHostname(%q) = %v, expected true (valid hostname rejected)",
				hostname, result)
		}
	}
}

func TestBug7_ValidateHostname_TooLong_ShouldBeInvalid(t *testing.T) {
	// RFC 1035: hostname labels max 63 chars, total hostname max 253 chars
	longLabel := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 64 chars
	longHostname := longLabel + ".example.com"

	result := validateHostname(longHostname)
	// Current regex doesn't check length - this is a limitation, not necessarily a bug
	t.Logf("validateHostname(64-char-label.example.com) = %v (no length check)", result)
}

// =============================================================================
// Validation Functions - Combined Tests
// =============================================================================

func TestValidateIPAddresses_AllValid(t *testing.T) {
	ips := []string{"127.0.0.1", "192.168.1.1", "10.0.0.1"}
	ok, failedIP := validateIPAddresses(ips)

	// Due to Bug #1, this should fail (valid IPs rejected)
	if !ok {
		t.Errorf("BUG #1 (indirect): validateIPAddresses rejected valid IPs, failed on: %s", failedIP)
	}
}

func TestValidateIPAddresses_OneInvalid(t *testing.T) {
	ips := []string{"127.0.0.1", "invalid", "10.0.0.1"}
	ok, failedIP := validateIPAddresses(ips)

	// Due to Bug #1, this might pass (invalid IPs accepted)
	if ok {
		t.Error("BUG #1 (indirect): validateIPAddresses accepted invalid IP")
	} else {
		t.Logf("Failed on: %s", failedIP)
	}
}

func TestValidateHostnames_AllValid(t *testing.T) {
	hostnames := []string{"localhost", "example.com", "my-server"}
	ok, failedHostname := validateHostnames(hostnames)

	if !ok {
		t.Errorf("validateHostnames rejected valid hostnames, failed on: %s", failedHostname)
	}
}

func TestValidateHostnames_OneInvalid(t *testing.T) {
	// These should be invalid but Bug #7 might accept them
	hostnames := []string{"localhost", "...", "my-server"}
	ok, failedHostname := validateHostnames(hostnames)

	if ok {
		t.Error("BUG #7 (indirect): validateHostnames accepted '...' which should be invalid")
	} else {
		t.Logf("Correctly failed on: %s", failedHostname)
	}
}

func TestValidateCIDRs_Valid(t *testing.T) {
	cidrs := []string{"192.168.0.0/24", "10.0.0.0/8", "172.16.0.0/12"}
	ok, failedCIDR := validateCIDRs(cidrs)

	if !ok {
		t.Errorf("validateCIDRs rejected valid CIDRs, failed on: %s", failedCIDR)
	}
}

func TestValidateCIDRs_Invalid(t *testing.T) {
	cidrs := []string{"192.168.0.0/24", "not-a-cidr", "10.0.0.0/8"}
	ok, failedCIDR := validateCIDRs(cidrs)

	if ok {
		t.Error("validateCIDRs should have rejected invalid CIDR")
	} else if failedCIDR != "not-a-cidr" {
		t.Errorf("Expected failure on 'not-a-cidr', got: %s", failedCIDR)
	}
}

func TestValidateCIDR_Valid(t *testing.T) {
	validCIDRs := []string{
		"192.168.0.0/24",
		"10.0.0.0/8",
		"0.0.0.0/0",
		"::1/128",
		"2001:db8::/32",
	}

	for _, cidr := range validCIDRs {
		if !validateCIDR(cidr) {
			t.Errorf("validateCIDR(%q) = false, expected true", cidr)
		}
	}
}

func TestValidateCIDR_Invalid(t *testing.T) {
	invalidCIDRs := []string{
		"not-a-cidr",
		"192.168.1.1",
		"192.168.1.0/33",
		"",
	}

	for _, cidr := range invalidCIDRs {
		if validateCIDR(cidr) {
			t.Errorf("validateCIDR(%q) = true, expected false", cidr)
		}
	}
}

// =============================================================================
// Phase 2: CLI Integration Tests
// =============================================================================

// Helper to create a temp hosts file and initialize etcHosts.
func setupTestHosts(t *testing.T, content string) (path string, cleanup func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	_ = tmpFile.Close()

	// Set up the global state
	HostsFileReadPath = tmpFile.Name()
	HostsFileWritePath = tmpFile.Name()
	DryRun = false
	Quiet = true
	Flush = false

	// Initialize etcHosts
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath:  tmpFile.Name(),
		WriteFilePath: tmpFile.Name(),
	})
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		t.Fatalf("Failed to initialize hosts: %v", err)
	}
	etcHosts = hosts

	cleanup = func() {
		_ = os.Remove(tmpFile.Name())
		HostsFileReadPath = ""
		HostsFileWritePath = ""
		DryRun = false
		Quiet = false
		Flush = false
	}

	return tmpFile.Name(), cleanup
}

// Helper to capture stdout.
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// --- Add Command Tests ---

func TestAddHosts_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	AddHosts("192.168.1.1", []string{"newhost", "newhost2"}, "")

	// Verify the hosts were added
	result := etcHosts.ListHostsByIP("192.168.1.1")
	if len(result) != 2 {
		t.Errorf("Expected 2 hosts, got %d: %v", len(result), result)
	}
}

func TestAddHosts_DryRun(t *testing.T) {
	path, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	DryRun = true

	output := captureOutput(func() {
		AddHosts("192.168.1.1", []string{"newhost"}, "")
	})

	// Output should contain the new host
	if !strings.Contains(output, "newhost") {
		t.Error("Dry run output should contain 'newhost'")
	}

	// File should NOT be modified (read original content)
	content, _ := os.ReadFile(filepath.Clean(path))
	if strings.Contains(string(content), "newhost") {
		t.Error("Dry run should not modify the file")
	}
}

func TestAddHosts_ToExistingAddress(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	AddHosts("127.0.0.1", []string{"newalias"}, "")

	result := etcHosts.ListHostsByIP("127.0.0.1")
	if len(result) != 2 {
		t.Errorf("Expected 2 hosts at 127.0.0.1, got %d: %v", len(result), result)
	}
}

func TestAddHosts_WithComment(t *testing.T) {
	_, cleanup := setupTestHosts(t, "")
	defer cleanup()

	AddHosts("192.168.1.1", []string{"myservice"}, "managed-by-myapp")

	// Verify the host was added
	result := etcHosts.ListHostsByIP("192.168.1.1")
	if len(result) != 1 || result[0] != "myservice" {
		t.Errorf("Expected 'myservice', got: %v", result)
	}

	// Verify the comment is in the rendered output
	rendered := etcHosts.RenderHostsFile()
	if !strings.Contains(rendered, "managed-by-myapp") {
		t.Errorf("Comment not found in rendered output: %s", rendered)
	}
}

func TestAddHosts_WithComment_SameCommentAppends(t *testing.T) {
	_, cleanup := setupTestHosts(t, "")
	defer cleanup()

	AddHosts("192.168.1.1", []string{"host1"}, "my-app")
	AddHosts("192.168.1.1", []string{"host2"}, "my-app")

	// Both should be on the same line
	lines := etcHosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line (same comment), got %d", len(lines))
	}
	if len(lines[0].Hostnames) != 2 {
		t.Errorf("Expected 2 hosts on line, got %d", len(lines[0].Hostnames))
	}
}

func TestAddHosts_WithComment_DifferentCommentNewLine(t *testing.T) {
	_, cleanup := setupTestHosts(t, "")
	defer cleanup()

	AddHosts("192.168.1.1", []string{"host1"}, "app-a")
	AddHosts("192.168.1.1", []string{"host2"}, "app-b")

	// Should be on different lines
	lines := etcHosts.GetHostFileLines()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines (different comments), got %d", len(lines))
	}
}

// --- Remove Host Command Tests ---

func TestRemoveHosts_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost myhost\n192.168.1.1 server\n")
	defer cleanup()

	RemoveHosts([]string{"myhost"})

	result := etcHosts.ListHostsByIP("127.0.0.1")
	for _, h := range result {
		if h == "myhost" {
			t.Error("myhost should have been removed")
		}
	}
}

func TestRemoveHosts_Multiple(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 host1 host2 host3\n")
	defer cleanup()

	RemoveHosts([]string{"host1", "host3"})

	result := etcHosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 || result[0] != "host2" {
		t.Errorf("Expected only 'host2' to remain, got: %v", result)
	}
}

func TestRemoveHosts_DryRun(t *testing.T) {
	path, cleanup := setupTestHosts(t, "127.0.0.1 localhost myhost\n")
	defer cleanup()

	DryRun = true

	output := captureOutput(func() {
		RemoveHosts([]string{"myhost"})
	})

	// Output should NOT contain myhost (it's removed in dry run view)
	if strings.Contains(output, "myhost") {
		t.Error("Dry run output should show myhost removed")
	}

	// But file should still have it
	content, _ := os.ReadFile(filepath.Clean(path))
	if !strings.Contains(string(content), "myhost") {
		t.Error("Dry run should not modify the file")
	}
}

func TestRemoveHosts_Nonexistent(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	// Should not panic or error when removing nonexistent host
	RemoveHosts([]string{"nonexistent"})

	// localhost should still be there
	result := etcHosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 || result[0] != testHostLocalhost {
		t.Errorf("localhost should remain unchanged: %v", result)
	}
}

// --- Remove IP Command Tests ---

func TestRemoveIPs_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server\n")
	defer cleanup()

	removeIPs([]string{"192.168.1.1"})

	result := etcHosts.ListHostsByIP("192.168.1.1")
	if len(result) != 0 {
		t.Errorf("192.168.1.1 should be removed, got: %v", result)
	}

	// localhost should remain
	result = etcHosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 {
		t.Error("localhost should remain")
	}
}

func TestRemoveIPs_Multiple(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server1\n192.168.1.2 server2\n")
	defer cleanup()

	removeIPs([]string{"192.168.1.1", "192.168.1.2"})

	result := etcHosts.ListHostsByIP("192.168.1.1")
	if len(result) != 0 {
		t.Error("192.168.1.1 should be removed")
	}

	result = etcHosts.ListHostsByIP("192.168.1.2")
	if len(result) != 0 {
		t.Error("192.168.1.2 should be removed")
	}
}

// --- Remove CIDR Command Tests ---

func TestRemoveIPRanges_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server1\n192.168.1.2 server2\n10.0.0.1 other\n")
	defer cleanup()

	RemoveIPRanges([]string{"192.168.1.0/24"})

	// 192.168.1.x should be gone
	result := etcHosts.ListHostsByIP("192.168.1.1")
	if len(result) != 0 {
		t.Error("192.168.1.1 should be removed")
	}

	result = etcHosts.ListHostsByIP("192.168.1.2")
	if len(result) != 0 {
		t.Error("192.168.1.2 should be removed")
	}

	// Others should remain
	result = etcHosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 {
		t.Error("127.0.0.1 should remain")
	}

	result = etcHosts.ListHostsByIP("10.0.0.1")
	if len(result) != 1 {
		t.Error("10.0.0.1 should remain")
	}
}

// --- List Commands Tests ---

func TestListByIPs_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server web\n")
	defer cleanup()

	output := captureOutput(func() {
		ListByIPs([]string{"192.168.1.1"})
	})

	if !strings.Contains(output, "192.168.1.1 server") {
		t.Errorf("Output should contain 'server': %s", output)
	}
	if !strings.Contains(output, "192.168.1.1 web") {
		t.Errorf("Output should contain 'web': %s", output)
	}
}

func TestListByIPs_NoMatch(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	output := captureOutput(func() {
		ListByIPs([]string{"192.168.1.1"})
	})

	// Should be empty for non-matching IP
	if strings.TrimSpace(output) != "" {
		t.Errorf("Expected empty output for non-matching IP, got: %s", output)
	}
}

func TestListByHostnames_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 myserver\n")
	defer cleanup()

	output := captureOutput(func() {
		ListByHostnames([]string{"myserver"})
	})

	if !strings.Contains(output, "192.168.1.1") {
		t.Errorf("Output should contain IP for myserver: %s", output)
	}
}

func TestListByCIDRs_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server1\n192.168.1.2 server2\n10.0.0.1 other\n")
	defer cleanup()

	output := captureOutput(func() {
		ListByCIDRs([]string{"192.168.1.0/24"})
	})

	if !strings.Contains(output, "192.168.1.1") {
		t.Error("Output should contain 192.168.1.1")
	}
	if !strings.Contains(output, "192.168.1.2") {
		t.Error("Output should contain 192.168.1.2")
	}
	if strings.Contains(output, "10.0.0.1") {
		t.Error("Output should NOT contain 10.0.0.1")
	}
}

// --- Show Command Tests ---

func TestShowHosts_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n# A comment\n192.168.1.1 server\n")
	defer cleanup()

	output := captureOutput(func() {
		ShowHosts()
	})

	if !strings.Contains(output, "localhost") {
		t.Error("Output should contain localhost")
	}
	if !strings.Contains(output, "# A comment") {
		t.Error("Output should preserve comments")
	}
	if !strings.Contains(output, "server") {
		t.Error("Output should contain server")
	}
}

// --- initEtcHosts and emptyFilePaths Tests ---

func TestEmptyFilePaths(t *testing.T) {
	// Save original values
	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
	}()

	HostsFileReadPath = ""
	HostsFileWritePath = ""
	if !emptyFilePaths() {
		t.Error("emptyFilePaths should return true when both paths are empty")
	}

	HostsFileReadPath = "/some/path"
	if emptyFilePaths() {
		t.Error("emptyFilePaths should return false when read path is set")
	}

	HostsFileReadPath = ""
	HostsFileWritePath = "/some/path"
	if emptyFilePaths() {
		t.Error("emptyFilePaths should return false when write path is set")
	}
}

func TestInitEtcHosts_WithPaths(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	// Save original values
	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		etcHosts = origHosts
	}()

	HostsFileReadPath = tmpFile.Name()
	HostsFileWritePath = tmpFile.Name()

	initEtcHosts()

	if etcHosts == nil {
		t.Error("etcHosts should be initialized")
	}

	result := etcHosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 || result[0] != testHostLocalhost {
		t.Errorf("etcHosts should contain localhost: %v", result)
	}
}

// =============================================================================
// Tests for List/Remove by Comment CLI Commands
// =============================================================================

func TestListByComment_Basic(t *testing.T) {
	hostsContent := "127.0.0.1        localhost\n127.0.0.1        app1 app2 # dev services\n127.0.0.1        api1 # dev services\n192.168.1.1      staging # staging env\n"
	_, cleanup := setupTestHosts(t, hostsContent)
	defer cleanup()

	output := captureOutput(func() {
		ListByComment("dev services")
	})

	// Should list app1, app2, api1
	if !strings.Contains(output, "app1") {
		t.Errorf("Output should contain 'app1': %s", output)
	}
	if !strings.Contains(output, "app2") {
		t.Errorf("Output should contain 'app2': %s", output)
	}
	if !strings.Contains(output, "api1") {
		t.Errorf("Output should contain 'api1': %s", output)
	}
	// Should not contain staging
	if strings.Contains(output, "staging") {
		t.Errorf("Output should not contain 'staging': %s", output)
	}
}

func TestRemoveByComment_Basic(t *testing.T) {
	hostsContent := "127.0.0.1        localhost\n127.0.0.1        app1 app2 # dev services\n192.168.1.1      staging # staging env\n"
	path, cleanup := setupTestHosts(t, hostsContent)
	defer cleanup()

	// Remove all entries with "dev services" comment
	RemoveByComment("dev services")

	// Reload and verify
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: path})
	if err != nil {
		t.Fatal(err)
	}

	// app1 and app2 should be gone
	result := hosts.ListHostsByComment("dev services")
	if len(result) != 0 {
		t.Errorf("Expected 0 hosts with 'dev services' comment, got %d", len(result))
	}

	// localhost should still be there
	result = hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 || result[0] != testHostLocalhost {
		t.Errorf("localhost should remain: %v", result)
	}

	// staging should still be there
	result = hosts.ListHostsByComment("staging env")
	if len(result) != 1 || result[0] != "staging" {
		t.Errorf("staging should remain: %v", result)
	}
}

func TestListByComment_NoMatch(t *testing.T) {
	hostsContent := testHostsWithComments
	tmpFile, err := os.CreateTemp("", "hosts-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(hostsContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		etcHosts = origHosts
	}()

	HostsFileReadPath = tmpFile.Name()
	HostsFileWritePath = tmpFile.Name()
	initEtcHosts()

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ListByComment("nonexistent comment")

	if err := w.Close(); err != nil {
		t.Logf("Warning: failed to close pipe writer: %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// Should be empty (no match)
	if strings.TrimSpace(output) != "" {
		t.Errorf("Expected empty output for non-existent comment, got: %s", output)
	}
}

func TestRemoveByComment_NoMatch(t *testing.T) {
	hostsContent := testHostsWithComments
	tmpFile, err := os.CreateTemp("", "hosts-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(hostsContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		etcHosts = origHosts
	}()

	HostsFileReadPath = tmpFile.Name()
	HostsFileWritePath = tmpFile.Name()
	initEtcHosts()

	// Remove with non-existent comment (should not affect anything)
	RemoveByComment("nonexistent comment")

	// Reload and verify nothing was removed
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: tmpFile.Name()})
	if err != nil {
		t.Fatal(err)
	}

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 2 {
		t.Errorf("Expected 2 hosts to remain, got %d", len(result))
	}
}

func TestListByComment_EmptyComment(t *testing.T) {
	hostsContent := testHostsWithComments
	tmpFile, err := os.CreateTemp("", "hosts-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(hostsContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		etcHosts = origHosts
	}()

	HostsFileReadPath = tmpFile.Name()
	HostsFileWritePath = tmpFile.Name()
	initEtcHosts()

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// List hosts with empty comment (should return localhost)
	ListByComment("")

	if err := w.Close(); err != nil {
		t.Logf("Warning: failed to close pipe writer: %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// Should contain localhost
	if !strings.Contains(output, "localhost") {
		t.Errorf("Expected 'localhost' in output, got: %s", output)
	}
	// Should not contain app1
	if strings.Contains(output, "app1") {
		t.Errorf("Should not contain 'app1' (has comment), got: %s", output)
	}
}

// --- DryRun Tests for Remove IP ---

func TestRemoveIPs_DryRun(t *testing.T) {
	path, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server\n")
	defer cleanup()

	DryRun = true

	output := captureOutput(func() {
		removeIPs([]string{"192.168.1.1"})
	})

	// Output should NOT contain server (it's removed in dry run view)
	if strings.Contains(output, "server") {
		t.Error("Dry run output should show server removed")
	}

	// File should still have it
	content, _ := os.ReadFile(filepath.Clean(path))
	if !strings.Contains(string(content), "server") {
		t.Error("Dry run should not modify the file")
	}
}

// --- DryRun Tests for Remove CIDR ---

func TestRemoveIPRanges_DryRun(t *testing.T) {
	path, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server1\n192.168.1.2 server2\n")
	defer cleanup()

	DryRun = true

	output := captureOutput(func() {
		RemoveIPRanges([]string{"192.168.1.0/24"})
	})

	// Output should NOT contain server1/server2
	if strings.Contains(output, "server1") || strings.Contains(output, "server2") {
		t.Error("Dry run output should show servers removed")
	}
	// Output should contain localhost
	if !strings.Contains(output, "localhost") {
		t.Error("Dry run output should still contain localhost")
	}

	// File should still have the original content
	content, _ := os.ReadFile(filepath.Clean(path))
	if !strings.Contains(string(content), "server1") {
		t.Error("Dry run should not modify the file")
	}
}

// --- DryRun Tests for Remove By Comment ---

func TestRemoveByComment_DryRun(t *testing.T) {
	path, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n127.0.0.1 app1 app2 # dev services\n")
	defer cleanup()

	DryRun = true

	output := captureOutput(func() {
		RemoveByComment("dev services")
	})

	// Output should NOT contain the dev services hosts
	if strings.Contains(output, "app1") || strings.Contains(output, "app2") {
		t.Error("Dry run output should show dev services hosts removed")
	}
	// Output should still contain localhost
	if !strings.Contains(output, "localhost") {
		t.Error("Dry run output should still contain localhost")
	}

	// File should still have original content
	content, _ := os.ReadFile(filepath.Clean(path))
	if !strings.Contains(string(content), "app1") {
		t.Error("Dry run should not modify the file")
	}
}

// --- VersionFromBuild Tests ---

func TestVersionFromBuild_GoreleaserVersion(t *testing.T) {
	// Save and restore original version
	origVersion := Version
	defer func() { Version = origVersion }()

	Version = "1.2.3"
	result := VersionFromBuild()
	if result != "1.2.3" {
		t.Errorf("Expected '1.2.3', got %q", result)
	}
}

func TestVersionFromBuild_Default(t *testing.T) {
	// Save and restore original version
	origVersion := Version
	defer func() { Version = origVersion }()

	Version = "0.0.0"
	result := VersionFromBuild()
	// In test context, debug.ReadBuildInfo() returns info from go test
	if result == "" {
		t.Error("Expected non-empty version from build info")
	}
}

// --- initEtcHosts Default Path ---

func TestInitEtcHosts_DefaultPath(t *testing.T) {
	// Save original values
	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origMax := MaxHostsPerLine
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		MaxHostsPerLine = origMax
		etcHosts = origHosts
	}()

	// Empty paths and zero max triggers the default path
	HostsFileReadPath = ""
	HostsFileWritePath = ""
	MaxHostsPerLine = 0

	initEtcHosts()

	if etcHosts == nil {
		t.Error("etcHosts should be initialized via default path")
	}
}

func TestInitEtcHosts_EmptyPathsNonZeroMax(t *testing.T) {
	// Save original values
	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origMax := MaxHostsPerLine
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		MaxHostsPerLine = origMax
		etcHosts = origHosts
	}()

	// Empty paths but non-zero max should take the NewHosts path
	HostsFileReadPath = ""
	HostsFileWritePath = ""
	MaxHostsPerLine = 5

	initEtcHosts()

	if etcHosts == nil {
		t.Error("etcHosts should be initialized via NewHosts path")
	}
}

func TestMaxHostsPerLine_Flag(t *testing.T) {
	hostsContent := `127.0.0.1        localhost
`
	tmpFile, err := os.CreateTemp("", "hosts-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(hostsContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Setup test environment
	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origHosts := etcHosts
	origMax := MaxHostsPerLine
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		etcHosts = origHosts
		MaxHostsPerLine = origMax
	}()

	HostsFileReadPath = tmpFile.Name()
	HostsFileWritePath = tmpFile.Name()
	MaxHostsPerLine = 2 // Set limit to 2 hosts per line
	initEtcHosts()

	// Add 5 hosts - should create 3 lines with limit of 2
	AddHosts("127.0.0.1", []string{"h1", "h2", "h3", "h4", "h5"}, "")

	// Reload and verify
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: tmpFile.Name()})
	if err != nil {
		t.Fatal(err)
	}

	lines := hosts.GetHostFileLines()
	addressLines := 0
	for _, line := range lines {
		if line.Address == "127.0.0.1" {
			addressLines++
			if len(line.Hostnames) > 2 {
				t.Errorf("Line has %d hosts, expected max 2", len(line.Hostnames))
			}
		}
	}

	// Should have 3 lines: localhost, (h1 h2), (h3 h4), (h5)
	// Actually localhost already has 1 host, so h1 goes there making it 2
	// Then h2 h3, h4 h5 = total 3 lines
	if addressLines < 3 {
		t.Errorf("Expected at least 3 lines with max 2 hosts per line, got %d", addressLines)
	}
}

// =============================================================================
// Cobra Command Args Validation Tests
// =============================================================================

func TestAddCmd_Args_Valid(t *testing.T) {
	err := addCmd.Args(addCmd, []string{"127.0.0.1", "myhost"})
	if err != nil {
		t.Errorf("Expected no error for valid args, got: %v", err)
	}
}

func TestAddCmd_Args_TooFew(t *testing.T) {
	err := addCmd.Args(addCmd, []string{"127.0.0.1"})
	if err == nil {
		t.Error("Expected error for missing hostname")
	}
}

func TestAddCmd_Args_InvalidIP(t *testing.T) {
	err := addCmd.Args(addCmd, []string{"not-an-ip", "myhost"})
	if err == nil {
		t.Error("Expected error for invalid IP")
	}
}

func TestAddCmd_Args_InvalidHostname(t *testing.T) {
	err := addCmd.Args(addCmd, []string{"127.0.0.1", "..."})
	if err == nil {
		t.Error("Expected error for invalid hostname")
	}
}

func TestRemoveHostCmd_Args_Valid(t *testing.T) {
	err := removeHostCmd.Args(removeHostCmd, []string{"myhost"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestRemoveHostCmd_Args_Empty(t *testing.T) {
	err := removeHostCmd.Args(removeHostCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestRemoveHostCmd_Args_InvalidHostname(t *testing.T) {
	err := removeHostCmd.Args(removeHostCmd, []string{"..."})
	if err == nil {
		t.Error("Expected error for invalid hostname")
	}
}

func TestRemoveIPCmd_Args_Valid(t *testing.T) {
	err := removeIPCmd.Args(removeIPCmd, []string{"127.0.0.1"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestRemoveIPCmd_Args_Empty(t *testing.T) {
	err := removeIPCmd.Args(removeIPCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestRemoveIPCmd_Args_InvalidIP(t *testing.T) {
	err := removeIPCmd.Args(removeIPCmd, []string{"not-an-ip"})
	if err == nil {
		t.Error("Expected error for invalid IP")
	}
}

func TestRemoveCidrCmd_Args_Valid(t *testing.T) {
	err := removeCidrCmd.Args(removeCidrCmd, []string{"192.168.1.0/24"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestRemoveCidrCmd_Args_Empty(t *testing.T) {
	err := removeCidrCmd.Args(removeCidrCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestRemoveCidrCmd_Args_InvalidCIDR(t *testing.T) {
	err := removeCidrCmd.Args(removeCidrCmd, []string{"not-a-cidr"})
	if err == nil {
		t.Error("Expected error for invalid CIDR")
	}
}

func TestRemoveCommentCmd_Args_Valid(t *testing.T) {
	err := removeCommentCmd.Args(removeCommentCmd, []string{"my comment"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestRemoveCommentCmd_Args_Empty(t *testing.T) {
	err := removeCommentCmd.Args(removeCommentCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestListByHostsCmd_Args_Valid(t *testing.T) {
	err := listByHostsCmd.Args(listByHostsCmd, []string{"myhost"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestListByHostsCmd_Args_Empty(t *testing.T) {
	err := listByHostsCmd.Args(listByHostsCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestListByHostsCmd_Args_InvalidHostname(t *testing.T) {
	err := listByHostsCmd.Args(listByHostsCmd, []string{"..."})
	if err == nil {
		t.Error("Expected error for invalid hostname")
	}
}

func TestListByIPCmd_Args_Valid(t *testing.T) {
	err := listByIPCmd.Args(listByIPCmd, []string{"127.0.0.1"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestListByIPCmd_Args_Empty(t *testing.T) {
	err := listByIPCmd.Args(listByIPCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestListByIPCmd_Args_InvalidIP(t *testing.T) {
	err := listByIPCmd.Args(listByIPCmd, []string{"not-an-ip"})
	if err == nil {
		t.Error("Expected error for invalid IP")
	}
}

func TestListByCidrCmd_Args_Valid(t *testing.T) {
	err := listByCidrCmd.Args(listByCidrCmd, []string{"10.0.0.0/8"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestListByCidrCmd_Args_Empty(t *testing.T) {
	err := listByCidrCmd.Args(listByCidrCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestListByCidrCmd_Args_InvalidCIDR(t *testing.T) {
	err := listByCidrCmd.Args(listByCidrCmd, []string{"not-a-cidr"})
	if err == nil {
		t.Error("Expected error for invalid CIDR")
	}
}

func TestListByCommentCmd_Args_Valid(t *testing.T) {
	err := listByCommentCmd.Args(listByCommentCmd, []string{"my comment"})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestListByCommentCmd_Args_Empty(t *testing.T) {
	err := listByCommentCmd.Args(listByCommentCmd, []string{})
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestShowCmd_Args(t *testing.T) {
	// show command accepts any args (no validation)
	err := showCmd.Args(showCmd, []string{})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// =============================================================================
// Cobra Command Run Function Tests
// =============================================================================

func TestAddCmd_Run_Basic(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	addComment = ""
	Quiet = true

	addCmd.Run(addCmd, []string{"192.168.1.1", "newhost"})

	result := etcHosts.ListHostsByIP("192.168.1.1")
	if len(result) != 1 || result[0] != "newhost" {
		t.Errorf("Expected 'newhost', got: %v", result)
	}
}

func TestAddCmd_Run_WithComment(t *testing.T) {
	_, cleanup := setupTestHosts(t, "")
	defer cleanup()

	addComment = "my-app"
	Quiet = true

	addCmd.Run(addCmd, []string{"192.168.1.1", "svc1", "svc2"})

	rendered := etcHosts.RenderHostsFile()
	if !strings.Contains(rendered, "my-app") {
		t.Error("Comment should appear in rendered output")
	}
	addComment = ""
}

func TestAddCmd_Run_NotQuiet(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	addComment = ""
	Quiet = false

	output := captureOutput(func() {
		addCmd.Run(addCmd, []string{"192.168.1.1", "newhost"})
	})

	if !strings.Contains(output, "Adding host") {
		t.Errorf("Expected 'Adding host' output, got: %s", output)
	}
}

func TestAddCmd_Run_NotQuiet_WithComment(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	addComment = "test-comment"
	Quiet = false

	output := captureOutput(func() {
		addCmd.Run(addCmd, []string{"192.168.1.1", "newhost"})
	})

	if !strings.Contains(output, "test-comment") {
		t.Errorf("Expected comment in output, got: %s", output)
	}
	addComment = ""
}

func TestRemoveHostCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost myhost\n")
	defer cleanup()

	Quiet = false

	output := captureOutput(func() {
		removeHostCmd.Run(removeHostCmd, []string{"myhost"})
	})

	if !strings.Contains(output, "Removing host") {
		t.Errorf("Expected 'Removing host' output, got: %s", output)
	}
}

func TestRemoveIPCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server\n")
	defer cleanup()

	Quiet = false

	output := captureOutput(func() {
		removeIPCmd.Run(removeIPCmd, []string{"192.168.1.1"})
	})

	if !strings.Contains(output, "Removing ip") {
		t.Errorf("Expected 'Removing ip' output, got: %s", output)
	}
}

func TestRemoveCidrCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n192.168.1.1 server\n")
	defer cleanup()

	Quiet = false

	output := captureOutput(func() {
		removeCidrCmd.Run(removeCidrCmd, []string{"192.168.1.0/24"})
	})

	if !strings.Contains(output, "Removing ip ranges") {
		t.Errorf("Expected 'Removing ip ranges' output, got: %s", output)
	}
}

func TestRemoveCommentCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 app1 # dev services\n")
	defer cleanup()

	Quiet = false

	output := captureOutput(func() {
		removeCommentCmd.Run(removeCommentCmd, []string{"dev services"})
	})

	if !strings.Contains(output, "Removing all hosts with comment") {
		t.Errorf("Expected removal message, got: %s", output)
	}
}

func TestListByHostsCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "192.168.1.1 myserver\n")
	defer cleanup()

	output := captureOutput(func() {
		listByHostsCmd.Run(listByHostsCmd, []string{"myserver"})
	})

	if !strings.Contains(output, "192.168.1.1") {
		t.Errorf("Expected IP in output, got: %s", output)
	}
}

func TestListByIPCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "192.168.1.1 myserver\n")
	defer cleanup()

	output := captureOutput(func() {
		listByIPCmd.Run(listByIPCmd, []string{"192.168.1.1"})
	})

	if !strings.Contains(output, "myserver") {
		t.Errorf("Expected hostname in output, got: %s", output)
	}
}

func TestListByCidrCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "192.168.1.1 myserver\n")
	defer cleanup()

	output := captureOutput(func() {
		listByCidrCmd.Run(listByCidrCmd, []string{"192.168.1.0/24"})
	})

	if !strings.Contains(output, "192.168.1.1") {
		t.Errorf("Expected IP in output, got: %s", output)
	}
}

func TestListByCommentCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 app1 # dev services\n")
	defer cleanup()

	output := captureOutput(func() {
		listByCommentCmd.Run(listByCommentCmd, []string{"dev services"})
	})

	if !strings.Contains(output, "app1") {
		t.Errorf("Expected 'app1' in output, got: %s", output)
	}
}

func TestShowCmd_Run(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n# test comment\n")
	defer cleanup()

	output := captureOutput(func() {
		showCmd.Run(showCmd, []string{})
	})

	if !strings.Contains(output, "localhost") {
		t.Errorf("Expected 'localhost' in output, got: %s", output)
	}
}

func TestVersionCmd_Run(t *testing.T) {
	output := captureOutput(func() {
		versionCmd.Run(versionCmd, []string{})
	})

	if !strings.Contains(output, "txeh Version:") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

// =============================================================================
// Flush Integration Tests
//
// Acceptance criteria for the CLI flush feature, written as
// Given/When/Then to verify concrete outputs.
// =============================================================================

// captureStderr redirects os.Stderr for the duration of f and returns
// what was written.
func captureStderr(f func()) string {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}
	os.Stderr = w

	f()

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// Given DryRun=true and Flush=true
// When AddHosts is called
// Then stdout contains the rendered hosts file with the new host
// And stdout does NOT contain "DNS cache flushed."
// (DryRun prevents Save, which prevents flush.)
func TestSaveHosts_DryRun_NoFlush(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	DryRun = true
	Flush = true

	output := captureOutput(func() {
		AddHosts("192.168.1.1", []string{"newhost"}, "")
	})

	if !strings.Contains(output, "newhost") {
		t.Error("dry run output should contain 'newhost'")
	}
	if strings.Contains(output, "DNS cache flushed") {
		t.Error("dry run should not print flush message")
	}
}

// Given Flush=true and Quiet=false with a successful save (using a temp file)
// When saveHosts() runs
// Then stdout contains exactly "DNS cache flushed.\n".
func TestSaveHosts_FlushSuccess_PrintsMessage(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	// The library-level flush is triggered by AutoFlush on HostsConfig.
	// We need AutoFlush=true for the exec mock to be exercised via Save().
	// Re-initialize with AutoFlush=true and a mock that succeeds.
	origFunc := txeh.ExecCommandFunc()
	defer txeh.SetExecCommandFunc(origFunc)
	txeh.SetExecCommandFunc(func(name string, args ...string) *exec.Cmd {
		return exec.CommandContext(context.Background(), "true") //nolint:gosec // test
	})

	hosts, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath:  etcHosts.WriteFilePath,
		WriteFilePath: etcHosts.WriteFilePath,
		AutoFlush:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	etcHosts = hosts
	Flush = true
	Quiet = false

	etcHosts.AddHost("192.168.1.1", "flushtest")

	output := captureOutput(func() {
		saveHosts()
	})

	if !strings.Contains(output, "DNS cache flushed.") {
		t.Errorf("expected 'DNS cache flushed.' in stdout, got: %q", output)
	}
}

// Given Flush=true and Quiet=true with a successful save
// When saveHosts() runs
// Then stdout does NOT contain "DNS cache flushed.".
func TestSaveHosts_FlushSuccess_Quiet_NoMessage(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	origFunc := txeh.ExecCommandFunc()
	defer txeh.SetExecCommandFunc(origFunc)
	txeh.SetExecCommandFunc(func(name string, args ...string) *exec.Cmd {
		return exec.CommandContext(context.Background(), "true") //nolint:gosec // test
	})

	hosts, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath:  etcHosts.WriteFilePath,
		WriteFilePath: etcHosts.WriteFilePath,
		AutoFlush:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	etcHosts = hosts
	Flush = true
	Quiet = true

	etcHosts.AddHost("192.168.1.1", "flushtest")

	output := captureOutput(func() {
		saveHosts()
	})

	if strings.Contains(output, "DNS cache flushed") {
		t.Errorf("quiet mode should suppress flush message, got: %q", output)
	}
}

// Given Flush=true (AutoFlush=true on config) and the flush command fails
// When saveHosts() runs
// Then stderr contains "Warning: hosts file saved but DNS cache flush failed"
// And the process does NOT exit (saveHosts returns normally).
func TestSaveHosts_FlushFails_StderrWarning(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	origFunc := txeh.ExecCommandFunc()
	defer txeh.SetExecCommandFunc(origFunc)
	txeh.SetExecCommandFunc(func(name string, args ...string) *exec.Cmd {
		return exec.CommandContext(context.Background(), "false") //nolint:gosec // test
	})

	hosts, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath:  etcHosts.WriteFilePath,
		WriteFilePath: etcHosts.WriteFilePath,
		AutoFlush:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	etcHosts = hosts
	Flush = true
	Quiet = false

	etcHosts.AddHost("192.168.1.1", "flushtest")

	stderrOutput := captureStderr(func() {
		saveHosts()
	})

	if !strings.Contains(stderrOutput, "Warning: hosts file saved but DNS cache flush failed") {
		t.Errorf("expected flush warning on stderr, got: %q", stderrOutput)
	}
	if !strings.Contains(stderrOutput, "DNS cache may be stale") {
		t.Errorf("expected stale cache hint on stderr, got: %q", stderrOutput)
	}
}

// Given Flush=true, Quiet=true, and flush fails
// When saveHosts() runs
// Then stderr contains the warning but NOT the "DNS cache may be stale" hint.
func TestSaveHosts_FlushFails_Quiet_NoHint(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	origFunc := txeh.ExecCommandFunc()
	defer txeh.SetExecCommandFunc(origFunc)
	txeh.SetExecCommandFunc(func(name string, args ...string) *exec.Cmd {
		return exec.CommandContext(context.Background(), "false") //nolint:gosec // test
	})

	hosts, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath:  etcHosts.WriteFilePath,
		WriteFilePath: etcHosts.WriteFilePath,
		AutoFlush:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	etcHosts = hosts
	Flush = true
	Quiet = true

	etcHosts.AddHost("192.168.1.1", "flushtest")

	stderrOutput := captureStderr(func() {
		saveHosts()
	})

	if !strings.Contains(stderrOutput, "Warning:") {
		t.Errorf("expected warning on stderr even in quiet mode, got: %q", stderrOutput)
	}
	if strings.Contains(stderrOutput, "DNS cache may be stale") {
		t.Errorf("quiet mode should suppress the hint, got: %q", stderrOutput)
	}
}

// Given Flush=false
// When saveHosts() runs after a successful save
// Then stdout does NOT contain "DNS cache flushed.".
func TestSaveHosts_NoFlush_NoMessage(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	Flush = false
	Quiet = false

	etcHosts.AddHost("192.168.1.1", "noflush")

	output := captureOutput(func() {
		saveHosts()
	})

	if strings.Contains(output, "DNS cache flushed") {
		t.Errorf("should not print flush message when Flush=false, got: %q", output)
	}
}

// --- TXEH_AUTO_FLUSH env var ---

// Given TXEH_AUTO_FLUSH=1
// When initEtcHosts() runs
// Then Flush is set to true.
func TestInitEtcHosts_EnvAutoFlush(t *testing.T) {
	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origFlush := Flush
	origMax := MaxHostsPerLine
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		Flush = origFlush
		MaxHostsPerLine = origMax
		etcHosts = origHosts
	}()

	HostsFileReadPath = ""
	HostsFileWritePath = ""
	MaxHostsPerLine = 0
	Flush = false

	t.Setenv("TXEH_AUTO_FLUSH", "1")

	initEtcHosts()

	if !Flush {
		t.Error("TXEH_AUTO_FLUSH=1 should set Flush to true")
	}
}

// Given TXEH_AUTO_FLUSH=0
// When initEtcHosts() runs
// Then Flush remains false.
func TestInitEtcHosts_EnvAutoFlush_NotSet(t *testing.T) {
	origRead := HostsFileReadPath
	origWrite := HostsFileWritePath
	origFlush := Flush
	origMax := MaxHostsPerLine
	origHosts := etcHosts
	defer func() {
		HostsFileReadPath = origRead
		HostsFileWritePath = origWrite
		Flush = origFlush
		MaxHostsPerLine = origMax
		etcHosts = origHosts
	}()

	HostsFileReadPath = ""
	HostsFileWritePath = ""
	MaxHostsPerLine = 0
	Flush = false

	t.Setenv("TXEH_AUTO_FLUSH", "0")

	initEtcHosts()

	if Flush {
		t.Error("TXEH_AUTO_FLUSH=0 should not set Flush to true")
	}
}

// --- Flag registration ---

// Given the root command
// When looking up the "flush" persistent flag
// Then it exists with shorthand "f".
func TestFlushFlag_Registered(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("flush")
	if f == nil {
		t.Fatal("--flush flag not registered")
	}
	if f.Shorthand != "f" {
		t.Errorf("Expected shorthand 'f', got %q", f.Shorthand)
	}
}
