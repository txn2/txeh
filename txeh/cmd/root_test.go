package cmd

import (
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/txn2/txeh"
)

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

// Helper to create a temp hosts file and initialize etcHosts
func setupTestHosts(t *testing.T, content string) (string, func()) {
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

	cleanup := func() {
		_ = os.Remove(tmpFile.Name())
		HostsFileReadPath = ""
		HostsFileWritePath = ""
		DryRun = false
		Quiet = false
	}

	return tmpFile.Name(), cleanup
}

// Helper to capture stdout
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

	AddHosts("192.168.1.1", []string{"newhost", "newhost2"})

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
		AddHosts("192.168.1.1", []string{"newhost"})
	})

	// Output should contain the new host
	if !strings.Contains(output, "newhost") {
		t.Error("Dry run output should contain 'newhost'")
	}

	// File should NOT be modified (read original content)
	content, _ := os.ReadFile(path)
	if strings.Contains(string(content), "newhost") {
		t.Error("Dry run should not modify the file")
	}
}

func TestAddHosts_ToExistingAddress(t *testing.T) {
	_, cleanup := setupTestHosts(t, "127.0.0.1 localhost\n")
	defer cleanup()

	AddHosts("127.0.0.1", []string{"newalias"})

	result := etcHosts.ListHostsByIP("127.0.0.1")
	if len(result) != 2 {
		t.Errorf("Expected 2 hosts at 127.0.0.1, got %d: %v", len(result), result)
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
	content, _ := os.ReadFile(path)
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
	if len(result) != 1 || result[0] != "localhost" {
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
	if len(result) != 1 || result[0] != "localhost" {
		t.Errorf("etcHosts should contain localhost: %v", result)
	}
}
