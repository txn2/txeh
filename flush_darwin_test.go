//go:build darwin

package txeh

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

// --- Platform flush command verification (darwin) ---

// Given execCommandFunc is mocked to record calls and succeed
// When FlushDNSCache() is called on darwin
// Then exactly 2 commands are invoked:
//
//	1st: "dscacheutil" with args ["-flushcache"]
//	2nd: "killall" with args ["-HUP", "mDNSResponder"]
//
// And no error is returned.
func TestFlushDNSCache_Darwin_CommandsAndArgs(t *testing.T) {
	origFunc := ExecCommandFunc()
	defer SetExecCommandFunc(origFunc)

	fn, calls, mu := mockExec("true")
	SetExecCommandFunc(fn)

	err := FlushDNSCache()
	if err != nil {
		t.Fatalf("FlushDNSCache() returned error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(*calls) != 2 {
		t.Fatalf("expected 2 exec calls, got %d: %+v", len(*calls), *calls)
	}

	// First call: dscacheutil -flushcache
	first := (*calls)[0]
	if first.Name != "dscacheutil" {
		t.Errorf("call[0].Name = %q, want %q", first.Name, "dscacheutil")
	}
	wantArgs := []string{"-flushcache"}
	if strings.Join(first.Args, " ") != strings.Join(wantArgs, " ") {
		t.Errorf("call[0].Args = %v, want %v", first.Args, wantArgs)
	}

	// Second call: killall -HUP mDNSResponder
	second := (*calls)[1]
	if second.Name != "killall" {
		t.Errorf("call[1].Name = %q, want %q", second.Name, "killall")
	}
	wantArgs2 := []string{"-HUP", "mDNSResponder"}
	if strings.Join(second.Args, " ") != strings.Join(wantArgs2, " ") {
		t.Errorf("call[1].Args = %v, want %v", second.Args, wantArgs2)
	}
}

// Given execCommandFunc is mocked to fail (exit 1)
// When FlushDNSCache() is called on darwin
// Then a *FlushError is returned with Platform="darwin"
// And Command="dscacheutil -flushcache"
// And only 1 command was attempted (killall is not reached).
func TestFlushDNSCache_Darwin_FirstCommandFails(t *testing.T) {
	origFunc := ExecCommandFunc()
	defer SetExecCommandFunc(origFunc)

	fn, calls, mu := mockExec("false")
	SetExecCommandFunc(fn)

	err := FlushDNSCache()
	if err == nil {
		t.Fatal("expected error when dscacheutil fails")
	}

	var fe *FlushError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *FlushError, got %T: %v", err, err)
	}

	if fe.Platform != "darwin" {
		t.Errorf("Platform = %q, want %q", fe.Platform, "darwin")
	}
	if fe.Command != "dscacheutil -flushcache" {
		t.Errorf("Command = %q, want %q", fe.Command, "dscacheutil -flushcache")
	}

	mu.Lock()
	defer mu.Unlock()

	// Only dscacheutil should have been attempted. killall should not run.
	if len(*calls) != 1 {
		t.Errorf("expected 1 exec call (dscacheutil only), got %d: %+v", len(*calls), *calls)
	}
}

// Given execCommandFunc where dscacheutil succeeds but killall fails
// When FlushDNSCache() is called
// Then no error is returned (killall failure is non-fatal).
func TestFlushDNSCache_Darwin_KillallFails_NonFatal(t *testing.T) {
	origFunc := ExecCommandFunc()
	defer SetExecCommandFunc(origFunc)

	callCount := 0
	var mu sync.Mutex
	SetExecCommandFunc(func(name string, args ...string) *exec.Cmd {
		mu.Lock()
		callCount++
		n := callCount
		mu.Unlock()

		if n == 1 {
			// dscacheutil succeeds
			return exec.CommandContext(context.Background(), "true") //nolint:gosec // test
		}
		// killall fails
		return exec.CommandContext(context.Background(), "false") //nolint:gosec // test
	})

	err := FlushDNSCache()
	if err != nil {
		t.Errorf("killall failure should be non-fatal, got error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if callCount != 2 {
		t.Errorf("expected 2 exec calls, got %d", callCount)
	}
}
