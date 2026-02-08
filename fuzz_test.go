package txeh

import (
	"strings"
	"testing"
)

const testBaseHosts = "127.0.0.1 existing\n"

// validateParsedHostLine checks that a parsed HostFileLine has valid type and
// consistent fields. ADDRESS lines must have a non-empty address and at least
// one hostname.
func validateParsedHostLine(t *testing.T, line HostFileLine) {
	t.Helper()
	switch line.LineType {
	case UNKNOWN, EMPTY, COMMENT, ADDRESS:
		// valid
	default:
		t.Errorf("invalid LineType %d", line.LineType)
	}
	if line.LineType == ADDRESS {
		if line.Address == "" {
			t.Error("ADDRESS line has empty Address")
		}
		if len(line.Hostnames) == 0 {
			t.Error("ADDRESS line has no hostnames")
		}
	}
}

// FuzzParseHostsFromString tests the hosts file parser with arbitrary input.
// This is the primary attack surface since it handles untrusted file content.
func FuzzParseHostsFromString(f *testing.F) {
	// Seed corpus with representative hosts file content
	f.Add("127.0.0.1 localhost")
	f.Add("127.0.0.1 localhost\n::1 localhost")
	f.Add("# comment line\n127.0.0.1 host1 host2 host3")
	f.Add("")
	f.Add("\n\n\n")
	f.Add("127.0.0.1 host # inline comment")
	f.Add("  \t  127.0.0.1  \t  host1  \t  host2  \t  ")
	f.Add("not-an-ip hostname")
	f.Add("127.0.0.1")
	f.Add("# only comments\n# more comments")
	f.Add("::1 ipv6host")
	f.Add("fe80::1%lo0 scoped-ipv6")
	f.Add("127.0.0.1 a b c d e f g h i j k l m n o p q r s t u v w x y z")
	f.Add("127.0.0.1 localhost\r\n::1 localhost\r\n")
	f.Add(strings.Repeat("127.0.0.1 host\n", 100))
	f.Add("127.0.0.1 host # comment # with # multiple # hashes")
	f.Add("#comment-no-space")
	f.Add("   ")
	f.Add("127.0.0.1 UPPERCASE lowercase MiXeD")

	f.Fuzz(func(t *testing.T, input string) {
		lines, err := ParseHostsFromString(input)
		if err != nil {
			return
		}

		// Every parsed line must have a valid type.
		for _, line := range lines {
			validateParsedHostLine(t, line)
		}
	})
}

// FuzzRoundTrip tests that parsing then rendering produces parseable output.
func FuzzRoundTrip(f *testing.F) {
	f.Add("127.0.0.1 localhost\n::1 localhost\n")
	f.Add("# comment\n127.0.0.1 host1 host2\n\n10.0.0.1 other\n")
	f.Add("127.0.0.1 host # with comment\n")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		hosts, err := NewHosts(&HostsConfig{RawText: &input})
		if err != nil {
			return
		}

		rendered := hosts.RenderHostsFile()

		// The rendered output must also be parseable
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("rendered output is not parseable: %v", err)
		}
	})
}

// FuzzAddHost tests adding hosts with arbitrary IP and hostname values.
func FuzzAddHost(f *testing.F) {
	f.Add("127.0.0.1", "localhost")
	f.Add("::1", "ipv6host")
	f.Add("10.0.0.1", "myhost.example.com")
	f.Add("not-an-ip", "hostname")
	f.Add("", "")
	f.Add("256.256.256.256", "host")
	f.Add("127.0.0.1", "UPPER")
	f.Add("127.0.0.1", "  spaces  ")
	f.Add("fe80::1", "v6host")

	f.Fuzz(func(t *testing.T, address, hostname string) {
		base := testBaseHosts
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		// Must not panic
		hosts.AddHost(address, hostname)

		// Must still render without error
		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("hosts file not parseable after AddHost(%q, %q): %v", address, hostname, err)
		}
	})
}

// FuzzAddHostWithComment tests adding hosts with arbitrary comments.
func FuzzAddHostWithComment(f *testing.F) {
	f.Add("127.0.0.1", "myhost", "added by app")
	f.Add("::1", "v6host", "ipv6 entry")
	f.Add("10.0.0.1", "host", "")
	f.Add("127.0.0.1", "host", "comment # with # hashes")
	f.Add("127.0.0.1", "host", "  leading and trailing spaces  ")

	f.Fuzz(func(t *testing.T, address, hostname, comment string) {
		base := testBaseHosts
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		hosts.AddHostWithComment(address, hostname, comment)

		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("hosts file not parseable after AddHostWithComment: %v", err)
		}
	})
}

// FuzzRemoveHost tests removing hosts with arbitrary hostnames.
func FuzzRemoveHost(f *testing.F) {
	f.Add("localhost")
	f.Add("nonexistent")
	f.Add("")
	f.Add("  spaces  ")
	f.Add("UPPERCASE")
	f.Add("host.with.dots")

	f.Fuzz(func(t *testing.T, hostname string) {
		base := "127.0.0.1 localhost host1 host2\n10.0.0.1 other\n::1 ip6host\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		// Must not panic
		hosts.RemoveHost(hostname)

		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("hosts file not parseable after RemoveHost(%q): %v", hostname, err)
		}
	})
}

// FuzzRemoveAddress tests removing addresses with arbitrary IP strings.
func FuzzRemoveAddress(f *testing.F) {
	f.Add("127.0.0.1")
	f.Add("10.0.0.1")
	f.Add("::1")
	f.Add("nonexistent")
	f.Add("")

	f.Fuzz(func(t *testing.T, address string) {
		base := "127.0.0.1 localhost\n10.0.0.1 other\n::1 ip6host\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		hosts.RemoveAddress(address)

		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("hosts file not parseable after RemoveAddress(%q): %v", address, err)
		}
	})
}

// FuzzRemoveCIDRs tests CIDR removal with arbitrary CIDR strings.
func FuzzRemoveCIDRs(f *testing.F) {
	f.Add("127.0.0.0/8")
	f.Add("10.0.0.0/24")
	f.Add("::1/128")
	f.Add("invalid-cidr")
	f.Add("")
	f.Add("192.168.1.0/32")
	f.Add("0.0.0.0/0")

	f.Fuzz(func(t *testing.T, cidr string) {
		base := "127.0.0.1 localhost\n10.0.0.1 other\n192.168.1.1 internal\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		// RemoveCIDRs may return an error for invalid CIDRs, that's expected
		_ = hosts.RemoveCIDRs([]string{cidr})

		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("hosts file not parseable after RemoveCIDRs(%q): %v", cidr, err)
		}
	})
}

// FuzzListOperations tests query operations with arbitrary input.
func FuzzListOperations(f *testing.F) {
	f.Add("localhost", "127.0.0.1", "10.0.0.0/8")
	f.Add("nonexistent", "256.256.256.256", "invalid")
	f.Add("", "", "")
	f.Add("UPPER", "::1", "::1/128")

	f.Fuzz(func(t *testing.T, hostname, address, cidr string) {
		base := "127.0.0.1 localhost host1\n10.0.0.1 other\n::1 ip6host\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		// None of these should panic
		hosts.ListHostsByIP(address)
		hosts.ListAddressesByHost(hostname, true)
		hosts.ListAddressesByHost(hostname, false)
		hosts.ListHostsByCIDR(cidr)
	})
}

// FuzzRemoveByComment tests comment-based removal with arbitrary strings.
func FuzzRemoveByComment(f *testing.F) {
	f.Add("added by app")
	f.Add("")
	f.Add("nonexistent comment")
	f.Add("  spaces  ")
	f.Add("comment # with # hashes")

	f.Fuzz(func(t *testing.T, comment string) {
		base := "127.0.0.1 host1 # added by app\n10.0.0.1 host2 # other comment\n127.0.0.2 host3\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		hosts.RemoveByComment(comment)

		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("hosts file not parseable after RemoveByComment(%q): %v", comment, err)
		}
	})
}

// FuzzListHostsByComment tests comment-based listing with arbitrary strings.
func FuzzListHostsByComment(f *testing.F) {
	f.Add("added by app")
	f.Add("")
	f.Add("nonexistent")

	f.Fuzz(func(t *testing.T, comment string) {
		base := "127.0.0.1 host1 # added by app\n10.0.0.1 host2\n"
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		// Must not panic
		hosts.ListHostsByComment(comment)
	})
}

// FuzzCombinedOperations tests sequences of add/remove/query with fuzzed data.
func FuzzCombinedOperations(f *testing.F) {
	f.Add("127.0.0.1", "host1", "10.0.0.1", "host2")
	f.Add("::1", "v6host", "fe80::1", "scoped")
	f.Add("", "", "", "")
	f.Add("not-ip", "host", "127.0.0.1", "  ")

	f.Fuzz(func(t *testing.T, addr1, host1, addr2, host2 string) {
		base := testBaseHosts
		hosts, err := NewHosts(&HostsConfig{RawText: &base})
		if err != nil {
			t.Fatalf("failed to create hosts: %v", err)
		}

		// Sequence of operations must not panic or corrupt state
		hosts.AddHost(addr1, host1)
		hosts.AddHost(addr2, host2)
		hosts.RemoveHost(host1)
		hosts.AddHost(addr1, host2)
		hosts.RemoveAddress(addr2)

		rendered := hosts.RenderHostsFile()
		_, err = ParseHostsFromString(rendered)
		if err != nil {
			t.Errorf("hosts file not parseable after combined operations: %v", err)
		}
	})
}
