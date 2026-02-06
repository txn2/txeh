package txeh

import (
	"strings"
	"testing"
	"testing/quick"

	"pgregory.net/rapid"
)

// =============================================================================
// Property-based tests using testing/quick (stdlib)
// These verify invariants that must hold for ALL valid inputs.
// =============================================================================

// TestProperty_ParseRenderRoundTrip_AlwaysParseable verifies that parsing any
// hosts content and rendering it always produces parseable output.
func TestProperty_ParseRenderRoundTrip_AlwaysParseable(t *testing.T) {
	t.Parallel()

	f := func(input string) bool {
		hosts, err := NewHosts(&HostsConfig{RawText: &input})
		if err != nil {
			return true // invalid input is fine to reject
		}
		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		return err == nil
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 1000}); err != nil {
		t.Error(err)
	}
}

// TestProperty_AddHost_ThenLookup_FindsIt verifies that adding a valid host
// to a valid address always makes it findable via lookup.
func TestProperty_AddHost_ThenLookup_FindsIt(t *testing.T) {
	t.Parallel()

	// Use a fixed set of valid IPs to test the property
	validIPs := []string{"10.0.0.1", "192.168.1.1", "172.16.0.1"}
	validHosts := []string{"host-a", "host-b", "host-c", "host-d", "host-e"}

	f := func(ipIdx, hostIdx uint) bool {
		ip := validIPs[ipIdx%uint(len(validIPs))]
		hostname := validHosts[hostIdx%uint(len(validHosts))]

		base := ""
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			return false
		}

		hosts.AddHost(ip, hostname)
		found := hosts.ListHostsByIP(ip)

		for _, h := range found {
			if h == hostname {
				return true
			}
		}
		return false
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Error(err)
	}
}

// TestProperty_RemoveHost_ThenLookup_NeverFindsIt verifies that removing a host
// guarantees it is no longer associated with any address.
func TestProperty_RemoveHost_ThenLookup_NeverFindsIt(t *testing.T) {
	t.Parallel()

	hostnames := []string{"alpha", "bravo", "charlie", "delta"}

	f := func(idx uint) bool {
		hostname := hostnames[idx%uint(len(hostnames))]

		base := "10.0.0.1 alpha bravo\n10.0.0.2 charlie delta\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			return false
		}

		hosts.RemoveHost(hostname)

		// Must not appear at any address
		lines := hosts.GetHostFileLines()
		for _, line := range lines {
			for _, h := range line.Hostnames {
				if h == hostname {
					return false
				}
			}
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Error(err)
	}
}

// TestProperty_RemoveAddress_RemovesAllHostnames verifies that removing an
// address removes every hostname that was on that address's lines.
func TestProperty_RemoveAddress_RemovesAllHostnames(t *testing.T) {
	t.Parallel()

	base := "10.0.0.1 host1 host2 host3\n10.0.0.2 host4\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &base})
	if err != nil {
		t.Fatal(err)
	}

	hosts.RemoveAddress("10.0.0.1")

	result := hosts.ListHostsByIP("10.0.0.1")
	if len(result) != 0 {
		t.Errorf("after RemoveAddress, expected 0 hosts at 10.0.0.1, got %d: %v", len(result), result)
	}

	// Other addresses must be unaffected
	other := hosts.ListHostsByIP("10.0.0.2")
	if len(other) != 1 || other[0] != "host4" {
		t.Errorf("RemoveAddress affected unrelated address, got: %v", other)
	}
}

// TestProperty_AddHost_InvalidIP_IsNoOp verifies that adding a host with an
// invalid IP address never modifies the hosts file.
func TestProperty_AddHost_InvalidIP_IsNoOp(t *testing.T) {
	t.Parallel()

	f := func(hostname string) bool {
		base := "127.0.0.1 existing\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			return true
		}

		before := hosts.RenderHostsFile()
		hosts.AddHost("not-a-valid-ip", hostname)
		after := hosts.RenderHostsFile()

		return before == after
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Error(err)
	}
}

// TestProperty_ParsedLineCount_MatchesInputLines verifies that parsing always
// produces the same number of lines as the input (minus trailing empty line).
func TestProperty_ParsedLineCount_MatchesInputLines(t *testing.T) {
	t.Parallel()

	f := func(input string) bool {
		lines, err := ParseHostsFromString(input)
		if err != nil {
			return true
		}

		normalized := strings.ReplaceAll(input, "\r\n", "\n")
		inputLines := strings.Split(normalized, "\n")
		// Parser trims trailing empty line if present
		if len(inputLines) > 0 && inputLines[len(inputLines)-1] == "" {
			inputLines = inputLines[:len(inputLines)-1]
		}

		return len(lines) == len(inputLines)
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 1000}); err != nil {
		t.Error(err)
	}
}

// =============================================================================
// Property-based tests using rapid
// These use controlled generators for more targeted invariant testing.
// =============================================================================

// validIPGen generates valid IPv4 addresses in private ranges.
func validIPGen() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		octet3 := rapid.IntRange(0, 255).Draw(t, "octet3")
		octet4 := rapid.IntRange(1, 254).Draw(t, "octet4")
		prefix := rapid.SampledFrom([]string{"10.0", "192.168", "172.16"}).Draw(t, "prefix")
		return prefix + "." + itoa(octet3) + "." + itoa(octet4)
	})
}

// hostnameGen generates valid-looking hostnames.
func hostnameGen() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		prefix := rapid.SampledFrom([]string{
			"app", "web", "api", "db", "cache", "svc", "node", "host",
		}).Draw(t, "prefix")
		suffix := rapid.SampledFrom([]string{
			".local", ".internal", ".test", ".example.com", "",
		}).Draw(t, "suffix")
		num := rapid.IntRange(0, 99).Draw(t, "num")
		return prefix + "-" + itoa(num) + suffix
	})
}

// itoa converts an int to a string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

// TestRapid_AddRemove_Symmetry verifies that adding then removing a host
// returns the hosts file to a state where that host is not present.
func TestRapid_AddRemove_Symmetry(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		ip := validIPGen().Draw(t, "ip")
		hostname := hostnameGen().Draw(t, "hostname")

		base := "127.0.0.1 localhost\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatal(err)
		}

		hosts.AddHost(ip, hostname)

		// Verify it was added
		found := hosts.ListHostsByIP(ip)
		hasHost := false
		for _, h := range found {
			if h == hostname {
				hasHost = true
				break
			}
		}
		if !hasHost {
			t.Fatalf("AddHost(%q, %q) did not add the host", ip, hostname)
		}

		// Remove and verify it's gone
		hosts.RemoveHost(hostname)
		foundAfter := hosts.ListHostsByIP(ip)
		for _, h := range foundAfter {
			if h == hostname {
				t.Fatalf("RemoveHost(%q) did not remove the host from %q", hostname, ip)
			}
		}
	})
}

// TestRapid_AddHost_Idempotent verifies that adding the same host twice
// to the same address doesn't create duplicates.
func TestRapid_AddHost_Idempotent(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		ip := validIPGen().Draw(t, "ip")
		hostname := hostnameGen().Draw(t, "hostname")

		base := ""
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatal(err)
		}

		hosts.AddHost(ip, hostname)
		hosts.AddHost(ip, hostname)

		found := hosts.ListHostsByIP(ip)
		count := 0
		for _, h := range found {
			if h == hostname {
				count++
			}
		}
		if count != 1 {
			t.Fatalf("AddHost called twice produced %d copies of %q (want 1)", count, hostname)
		}
	})
}

// TestRapid_RemoveHost_Idempotent verifies that removing a host that doesn't
// exist is a safe no-op that doesn't corrupt state.
func TestRapid_RemoveHost_Idempotent(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		hostname := hostnameGen().Draw(t, "hostname")

		base := "127.0.0.1 localhost\n10.0.0.1 existing\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatal(err)
		}

		before := hosts.RenderHostsFile()
		hosts.RemoveHost(hostname)
		hosts.RemoveHost(hostname) // second remove should be safe no-op
		after := hosts.RenderHostsFile()

		// If host wasn't in original, content should be unchanged
		if !strings.Contains(base, hostname) && before != after {
			t.Fatalf("removing non-existent host %q modified the file", hostname)
		}
	})
}

// TestRapid_BulkAddRemove_StateConsistency verifies that after a sequence of
// adds and removes, the rendered hosts file is always valid and consistent.
func TestRapid_BulkAddRemove_StateConsistency(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		numOps := rapid.IntRange(5, 30).Draw(t, "numOps")

		base := "127.0.0.1 localhost\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < numOps; i++ {
			ip := validIPGen().Draw(t, "ip")
			hostname := hostnameGen().Draw(t, "hostname")

			if rapid.Bool().Draw(t, "isAdd") {
				hosts.AddHost(ip, hostname)
			} else {
				hosts.RemoveHost(hostname)
			}
		}

		// After all operations, state must be consistent
		rendered := hosts.RenderHostsFile()

		// Must be parseable
		reparsed, err := ParseHostsFromString(rendered)
		if err != nil {
			t.Fatalf("rendered output not parseable after %d ops: %v", numOps, err)
		}

		// All ADDRESS lines must have non-empty address and hostnames
		for _, line := range reparsed {
			if line.LineType == ADDRESS {
				if line.Address == "" {
					t.Fatal("ADDRESS line has empty address after bulk operations")
				}
				if len(line.Hostnames) == 0 {
					t.Fatal("ADDRESS line has no hostnames after bulk operations")
				}
			}
		}
	})
}

// TestRapid_CommentOperations_Consistent verifies that adding hosts with comments
// and removing by comment produces consistent state.
func TestRapid_CommentOperations_Consistent(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		ip := validIPGen().Draw(t, "ip")
		hostname := hostnameGen().Draw(t, "hostname")
		comment := rapid.SampledFrom([]string{
			"kubefwd", "added by myapp", "test entry", "k8s-svc",
		}).Draw(t, "comment")

		base := "127.0.0.1 localhost\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatal(err)
		}

		// Add with comment
		hosts.AddHostWithComment(ip, hostname, comment)

		// Must be findable by comment
		byComment := hosts.ListHostsByComment(comment)
		found := false
		for _, h := range byComment {
			if h == hostname {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("AddHostWithComment(%q, %q, %q): host not found via ListHostsByComment", ip, hostname, comment)
		}

		// Remove by comment
		hosts.RemoveByComment(comment)

		// Must no longer be findable by comment
		byCommentAfter := hosts.ListHostsByComment(comment)
		if len(byCommentAfter) != 0 {
			t.Fatalf("RemoveByComment(%q) did not remove hosts: %v", comment, byCommentAfter)
		}
	})
}

// TestRapid_HostnameMoveBetweenAddresses verifies that adding a hostname to a
// new (non-localhost) address removes it from the old address.
func TestRapid_HostnameMoveBetweenAddresses(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		hostname := hostnameGen().Draw(t, "hostname")
		// Use non-localhost IPs so the move behavior applies
		ip1 := "10." + itoa(rapid.IntRange(0, 255).Draw(t, "o1")) + ".0.1"
		ip2 := "10." + itoa(rapid.IntRange(0, 255).Draw(t, "o2")) + ".0.2"

		base := ""
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatal(err)
		}

		// Add to first address
		hosts.AddHost(ip1, hostname)

		// Move to second address
		hosts.AddHost(ip2, hostname)

		// If IPs differ, hostname should only be at ip2
		if ip1 != ip2 {
			atIP1 := hosts.ListHostsByIP(ip1)
			for _, h := range atIP1 {
				if h == hostname {
					t.Fatalf("hostname %q should have moved from %s to %s, but still at old address", hostname, ip1, ip2)
				}
			}
		}

		atIP2 := hosts.ListHostsByIP(ip2)
		found := false
		for _, h := range atIP2 {
			if h == hostname {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("hostname %q not found at new address %s", hostname, ip2)
		}
	})
}

// TestRapid_RenderParse_PreservesAddressEntries verifies that rendering and
// re-parsing preserves all address-hostname mappings.
func TestRapid_RenderParse_PreservesAddressEntries(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		numEntries := rapid.IntRange(1, 10).Draw(t, "numEntries")

		base := ""
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatal(err)
		}

		type entry struct {
			ip       string
			hostname string
		}
		var entries []entry
		for i := 0; i < numEntries; i++ {
			ip := validIPGen().Draw(t, "ip")
			hostname := hostnameGen().Draw(t, "hostname")
			hosts.AddHost(ip, hostname)
			entries = append(entries, entry{ip, hostname})
		}

		// Render and re-parse
		rendered := hosts.RenderHostsFile()
		hosts2, err := NewHosts(&HostsConfig{RawText: &rendered})
		if err != nil {
			t.Fatalf("re-parse failed: %v", err)
		}

		// Every entry that was added should still be findable
		// (some may have been moved if the same hostname was added to multiple IPs)
		for _, e := range entries {
			results := hosts2.ListAddressesByHost(e.hostname, true)
			if len(results) == 0 {
				t.Fatalf("hostname %q lost after render/re-parse", e.hostname)
			}
		}
	})
}
