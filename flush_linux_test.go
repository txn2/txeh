//go:build linux

package txeh

import (
	"errors"
	"strings"
	"testing"
)

// --- Platform flush command verification (linux) ---

// Given execCommandFunc is mocked to succeed and resolvectl is on PATH
// When FlushDNSCache() is called on linux
// Then "resolvectl flush-caches" is invoked
// And no error is returned.
func TestFlushDNSCache_Linux_ResolvectlSuccess(t *testing.T) {
	origFunc := ExecCommandFunc()
	defer SetExecCommandFunc(origFunc)

	fn, calls, mu := mockExec("true")
	SetExecCommandFunc(fn)

	err := FlushDNSCache()
	// On CI (ubuntu-latest), resolvectl is present. If neither binary exists,
	// ErrNoResolver is returned, which is also valid behavior.
	if err != nil {
		var fe *FlushError
		if errors.As(err, &fe) && errors.Is(fe.Err, ErrNoResolver) {
			t.Skip("neither resolvectl nor systemd-resolve on PATH, skipping")
		}
		t.Fatalf("FlushDNSCache() returned error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(*calls) != 1 {
		t.Fatalf("expected 1 exec call, got %d: %+v", len(*calls), *calls)
	}

	first := (*calls)[0]
	// Should be resolvectl or systemd-resolve depending on what's installed.
	if first.Name != "resolvectl" && first.Name != "systemd-resolve" {
		t.Errorf("call[0].Name = %q, want resolvectl or systemd-resolve", first.Name)
	}
}

// Given execCommandFunc is mocked to fail (exit 1) and resolvectl is on PATH
// When FlushDNSCache() is called on linux
// Then a *FlushError is returned with Platform="linux".
func TestFlushDNSCache_Linux_CommandFails(t *testing.T) {
	origFunc := ExecCommandFunc()
	defer SetExecCommandFunc(origFunc)

	fn, _, _ := mockExec("false")
	SetExecCommandFunc(fn)

	err := FlushDNSCache()
	if err == nil {
		t.Fatal("expected error when flush command fails")
	}

	var fe *FlushError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *FlushError, got %T: %v", err, err)
	}

	if fe.Platform != "linux" {
		t.Errorf("Platform = %q, want %q", fe.Platform, "linux")
	}
}

// Given the ErrNoResolver sentinel
// When its message is examined
// Then it mentions systemd-resolved and explains that flushing may not be needed.
func TestErrNoResolver_Message(t *testing.T) {
	t.Parallel()

	msg := ErrNoResolver.Error()
	if !strings.Contains(msg, "systemd-resolved not detected") {
		t.Errorf("ErrNoResolver should mention systemd-resolved, got: %q", msg)
	}
	if !strings.Contains(msg, "immediately without flushing") {
		t.Errorf("ErrNoResolver should explain no flush needed, got: %q", msg)
	}
}
