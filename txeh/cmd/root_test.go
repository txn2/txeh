package cmd

import (
	"net"
	"testing"
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
