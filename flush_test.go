package txeh

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

// cmdRecord captures a single call to execCommandFunc.
type cmdRecord struct {
	Name string
	Args []string
}

// mockExec returns a mock execCommandFunc that records every call and
// runs exitCmd (e.g. "true" or "false") instead of the real binary.
// The returned slice is append-safe across calls.
func mockExec(exitCmd string) (func(string, ...string) *exec.Cmd, *[]cmdRecord, *sync.Mutex) {
	var mu sync.Mutex
	var calls []cmdRecord
	fn := func(name string, args ...string) *exec.Cmd {
		mu.Lock()
		calls = append(calls, cmdRecord{Name: name, Args: args})
		mu.Unlock()
		return exec.CommandContext(context.Background(), exitCmd) //nolint:gosec // test helper
	}
	return fn, &calls, &mu
}

// --- FlushError acceptance criteria ---

// Given a FlushError with Platform="darwin", Command="dscacheutil -flushcache", Err="exit status 1"
// When Error() is called
// Then the string is exactly "flush DNS cache on darwin (dscacheutil -flushcache): exit status 1".
func TestFlushError_Error_ExactFormat(t *testing.T) {
	t.Parallel()

	fe := &FlushError{
		Platform: "darwin",
		Command:  "dscacheutil -flushcache",
		Err:      fmt.Errorf("exit status 1"),
	}

	want := "flush DNS cache on darwin (dscacheutil -flushcache): exit status 1"
	got := fe.Error()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Given a FlushError with an empty Command field
// When Error() is called
// Then the parentheses are still present but empty.
func TestFlushError_Error_EmptyCommand(t *testing.T) {
	t.Parallel()

	fe := &FlushError{
		Platform: "freebsd",
		Command:  "",
		Err:      fmt.Errorf("unsupported platform"),
	}

	want := "flush DNS cache on freebsd (): unsupported platform"
	got := fe.Error()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Given a FlushError wrapping a sentinel error
// When checked with errors.Is
// Then the sentinel is found through Unwrap.
func TestFlushError_Unwrap_ErrorsIs(t *testing.T) {
	t.Parallel()

	sentinel := fmt.Errorf("inner cause")
	fe := &FlushError{Platform: "linux", Command: "resolvectl flush-caches", Err: sentinel}

	if !errors.Is(fe, sentinel) {
		t.Error("errors.Is should find the inner sentinel through Unwrap")
	}
}

// Given a *FlushError stored as an error interface
// When checked with errors.As
// Then it extracts the *FlushError with correct Platform and Command fields.
func TestFlushError_ErrorsAs_ExtractsFields(t *testing.T) {
	t.Parallel()

	var err error = &FlushError{
		Platform: "windows",
		Command:  "ipconfig /flushdns",
		Err:      fmt.Errorf("access denied"),
	}

	var fe *FlushError
	if !errors.As(err, &fe) {
		t.Fatal("errors.As should match *FlushError")
	}
	if fe.Platform != "windows" {
		t.Errorf("Platform = %q, want %q", fe.Platform, "windows")
	}
	if fe.Command != "ipconfig /flushdns" {
		t.Errorf("Command = %q, want %q", fe.Command, "ipconfig /flushdns")
	}
}

// --- joinArgs acceptance criteria ---

// Given various command + args combinations
// When joinArgs is called
// Then the output matches the exact expected string.
func TestJoinArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cmd  string
		args []string
		want string
	}{
		{name: "no args", cmd: "dscacheutil", args: nil, want: "dscacheutil"},
		{name: "single arg", cmd: "dscacheutil", args: []string{"-flushcache"}, want: "dscacheutil -flushcache"},
		{name: "multiple args", cmd: "killall", args: []string{"-HUP", "mDNSResponder"}, want: "killall -HUP mDNSResponder"},
		{name: "empty args slice", cmd: "ipconfig", args: []string{}, want: "ipconfig"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := joinArgs(tc.cmd, tc.args)
			if got != tc.want {
				t.Errorf("joinArgs(%q, %v) = %q, want %q", tc.cmd, tc.args, got, tc.want)
			}
		})
	}
}

// --- SaveAs + AutoFlush integration ---

// Given a Hosts instance with AutoFlush=true and a mocked exec
// When Save() is called
// Then the file is written to disk
// And FlushDNSCache is invoked (execCommandFunc is called).
func TestSaveAs_AutoFlush_True_FileWrittenAndFlushCalled(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	fn, calls, mu := mockExec("true")
	execCommandFunc = fn

	tmpFile, err := os.CreateTemp("", "hosts_flush_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	hosts, err := NewHosts(&HostsConfig{
		ReadFilePath:  tmpFile.Name(),
		WriteFilePath: tmpFile.Name(),
		AutoFlush:     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	hosts.AddHost("192.168.1.1", "testhost")
	if err := hosts.Save(); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file was actually written.
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "testhost") {
		t.Error("file should contain 'testhost' after Save()")
	}

	// Verify flush was called.
	mu.Lock()
	defer mu.Unlock()
	if len(*calls) == 0 {
		t.Error("AutoFlush=true: expected at least 1 exec call, got 0")
	}
}

// Given a Hosts instance with AutoFlush=false and a mocked exec
// When Save() is called
// Then the file is written to disk
// And execCommandFunc is never called.
func TestSaveAs_AutoFlush_False_FileWrittenNoFlush(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	fn, calls, mu := mockExec("true")
	execCommandFunc = fn

	tmpFile, err := os.CreateTemp("", "hosts_flush_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	hosts, err := NewHosts(&HostsConfig{
		ReadFilePath:  tmpFile.Name(),
		WriteFilePath: tmpFile.Name(),
		AutoFlush:     false,
	})
	if err != nil {
		t.Fatal(err)
	}

	hosts.AddHost("192.168.1.1", "testhost")
	if err := hosts.Save(); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file was written.
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "testhost") {
		t.Error("file should contain 'testhost' after Save()")
	}

	// Verify flush was NOT called.
	mu.Lock()
	defer mu.Unlock()
	if len(*calls) != 0 {
		t.Errorf("AutoFlush=false: expected 0 exec calls, got %d", len(*calls))
	}
}

// Given AutoFlush=true and a mock that fails
// When Save() is called
// Then it returns a *FlushError (not a write error)
// And the file was still written successfully (flush happens after write).
func TestSaveAs_AutoFlush_FlushFails_FileStillWritten(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	fn, _, _ := mockExec("false")
	execCommandFunc = fn

	tmpFile, err := os.CreateTemp("", "hosts_flush_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	hosts, err := NewHosts(&HostsConfig{
		ReadFilePath:  tmpFile.Name(),
		WriteFilePath: tmpFile.Name(),
		AutoFlush:     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	hosts.AddHost("192.168.1.1", "testhost")
	err = hosts.Save()

	// Should get an error.
	if err == nil {
		t.Fatal("expected error from Save() when flush fails")
	}

	// Error should be *FlushError, not a write error.
	var fe *FlushError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *FlushError, got %T: %v", err, err)
	}

	// File should still have been written (flush happens after write).
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "testhost") {
		t.Error("file should contain 'testhost' even when flush fails (write happens first)")
	}
}

// Given a Hosts instance created with RawText (no file)
// When AutoFlush is set
// Then Save() still returns the "cannot call Save with RawText" error, not a flush error.
func TestSaveAs_AutoFlush_RawText_ReturnsWriteError(t *testing.T) {
	raw := "127.0.0.1 localhost\n"
	hosts, err := NewHosts(&HostsConfig{
		RawText:   &raw,
		AutoFlush: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = hosts.Save()
	if err == nil {
		t.Fatal("expected error from Save() with RawText")
	}

	// Should be a plain error, not a *FlushError.
	var fe *FlushError
	if errors.As(err, &fe) {
		t.Error("RawText Save error should not be *FlushError")
	}

	if !strings.Contains(err.Error(), "RawText") {
		t.Errorf("expected RawText error message, got: %v", err)
	}
}
