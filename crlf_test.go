package txeh

import (
	"strings"
	"testing"
)

// TestStripLineBreaks verifies the helper drops only line-structure-breaking
// control characters and leaves all other content untouched.
func TestStripLineBreaks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain text untouched", "added by kubefwd", "added by kubefwd"},
		{"strips lf", "a\nb", "ab"},
		{"strips cr", "a\rb", "ab"},
		{"strips crlf", "a\r\nb", "ab"},
		{"strips vertical tab", "a\vb", "ab"},
		{"strips form feed", "a\fb", "ab"},
		{"strips nul", "a\x00b", "ab"},
		{"keeps tab", "a\tb", "a\tb"},
		{"keeps space", "a b", "a b"},
		{"keeps hash", "a # b", "a # b"},
		{"empty string", "", ""},
		{"only breaks", "\r\n\v\f\x00", ""},
		{"crlf injected ssh key", "x\r\nssh-rsa AAAA attacker@evil", "xssh-rsa AAAA attacker@evil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := stripLineBreaks(tt.in); got != tt.want {
				t.Errorf("stripLineBreaks(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestAddHostWithComment_CRLF_ProducesSingleLine is the regression test for the
// reported CRLF injection. A comment carrying an embedded "\r\n" must not break
// a single host entry into multiple physical lines.
func TestAddHostWithComment_CRLF_ProducesSingleLine(t *testing.T) {
	t.Parallel()

	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("NewHosts: %v", err)
	}

	// The sudoers-style payload from the PoC.
	hosts.AddHostWithComment(testIPv4Localhost, "placeholder", "x\r\nattacker ALL=(ALL) NOPASSWD: ALL")

	rendered := hosts.RenderHostsFile()

	if strings.ContainsAny(rendered, "\r") {
		t.Errorf("rendered output contains carriage return: %q", rendered)
	}

	// Exactly one non-empty physical line should be produced.
	var nonEmpty []string
	for ln := range strings.SplitSeq(rendered, "\n") {
		if strings.TrimSpace(ln) != "" {
			nonEmpty = append(nonEmpty, ln)
		}
	}
	if len(nonEmpty) != 1 {
		t.Fatalf("expected 1 non-empty line, got %d: %#v", len(nonEmpty), nonEmpty)
	}

	want := "127.0.0.1        placeholder # xattacker ALL=(ALL) NOPASSWD: ALL"
	if nonEmpty[0] != want {
		t.Errorf("line = %q, want %q", nonEmpty[0], want)
	}
}

// TestAddHost_CRLF_InHostname_ProducesSingleLine verifies the hostname field is
// also neutralized, not just comments.
func TestAddHost_CRLF_InHostname_ProducesSingleLine(t *testing.T) {
	t.Parallel()

	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("NewHosts: %v", err)
	}

	hosts.AddHost(testIPv4Localhost, "good\r\nevil")

	rendered := hosts.RenderHostsFile()
	if strings.ContainsAny(rendered, "\r") || strings.Count(strings.TrimRight(rendered, "\n"), "\n") != 0 {
		t.Errorf("hostname line break not neutralized: %q", rendered)
	}
}

// TestLineFormatter_CRLF_Backstop verifies the render-time guard catches
// line-break characters even when fields are populated directly, bypassing
// addHostWithComment.
func TestLineFormatter_CRLF_Backstop(t *testing.T) {
	t.Parallel()

	hfl := HostFileLine{
		LineType:  ADDRESS,
		Address:   testIPv4Localhost,
		Hostnames: []string{"host\r\ninjected"},
		Comment:   "note\r\nmore",
	}

	got := lineFormatter(hfl)
	want := "127.0.0.1        hostinjected # notemore"
	if got != want {
		t.Errorf("lineFormatter = %q, want %q", got, want)
	}
}
