package txeh

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestFlushError_Error(t *testing.T) {
	t.Parallel()

	fe := &FlushError{
		Platform: "darwin",
		Command:  "dscacheutil -flushcache",
		Err:      fmt.Errorf("exit status 1"),
	}

	want := "flush DNS cache on darwin (dscacheutil -flushcache): exit status 1"
	got := fe.Error()
	if got != want {
		t.Errorf("FlushError.Error() = %q, want %q", got, want)
	}
}

func TestFlushError_Unwrap(t *testing.T) {
	t.Parallel()

	inner := fmt.Errorf("something broke")
	fe := &FlushError{
		Platform: "linux",
		Command:  "resolvectl flush-caches",
		Err:      inner,
	}

	if !errors.Is(fe, inner) {
		t.Error("errors.Is should find the inner error via Unwrap")
	}
}

func TestFlushError_ErrorsAs(t *testing.T) {
	t.Parallel()

	inner := fmt.Errorf("cmd failed")
	var err error = &FlushError{
		Platform: "windows",
		Command:  "ipconfig /flushdns",
		Err:      inner,
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

func TestFlushError_EmptyCommand(t *testing.T) {
	t.Parallel()

	fe := &FlushError{
		Platform: "freebsd",
		Command:  "",
		Err:      fmt.Errorf("unsupported platform"),
	}

	want := "flush DNS cache on freebsd (): unsupported platform"
	got := fe.Error()
	if got != want {
		t.Errorf("FlushError.Error() = %q, want %q", got, want)
	}
}

func TestJoinArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cmd  string
		args []string
		want string
	}{
		{
			name: "no args",
			cmd:  "dscacheutil",
			args: nil,
			want: "dscacheutil",
		},
		{
			name: "single arg",
			cmd:  "dscacheutil",
			args: []string{"-flushcache"},
			want: "dscacheutil -flushcache",
		},
		{
			name: "multiple args",
			cmd:  "killall",
			args: []string{"-HUP", "mDNSResponder"},
			want: "killall -HUP mDNSResponder",
		},
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

// TestFlushDNSCache_MockSuccess replaces execCommandFunc to simulate a
// successful flush without running real system commands.
func TestFlushDNSCache_MockSuccess(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	// Mock exec.Command to return a command that exits 0.
	execCommandFunc = func(name string, args ...string) *exec.Cmd {
		return exec.CommandContext(context.Background(), "true")
	}

	err := FlushDNSCache()
	if err != nil {
		t.Errorf("FlushDNSCache() with mock success returned error: %v", err)
	}
}

// TestFlushDNSCache_MockFailure replaces execCommandFunc to simulate a
// failed flush and verifies the returned error type.
func TestFlushDNSCache_MockFailure(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	// Mock exec.Command to return a command that exits 1.
	execCommandFunc = func(name string, args ...string) *exec.Cmd {
		return exec.CommandContext(context.Background(), "false")
	}

	err := FlushDNSCache()
	if err == nil {
		t.Fatal("FlushDNSCache() with mock failure should return error")
	}

	var fe *FlushError
	if !errors.As(err, &fe) {
		t.Errorf("Expected *FlushError, got %T: %v", err, err)
	}
}

// TestSaveAs_AutoFlush_True verifies that SaveAs calls FlushDNSCache when
// AutoFlush is true.
func TestSaveAs_AutoFlush_True(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	flushCalled := false
	execCommandFunc = func(name string, args ...string) *exec.Cmd {
		flushCalled = true
		return exec.CommandContext(context.Background(), "true")
	}

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
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	if !flushCalled {
		t.Error("AutoFlush=true should have triggered FlushDNSCache")
	}
}

// TestSaveAs_AutoFlush_False verifies that SaveAs does NOT call FlushDNSCache
// when AutoFlush is false (default).
func TestSaveAs_AutoFlush_False(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	flushCalled := false
	execCommandFunc = func(name string, args ...string) *exec.Cmd {
		flushCalled = true
		return exec.CommandContext(context.Background(), "true")
	}

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
	err = hosts.Save()
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	if flushCalled {
		t.Error("AutoFlush=false should not trigger FlushDNSCache")
	}
}

// TestSaveAs_AutoFlush_Error verifies that a flush failure is returned
// as a *FlushError from Save().
func TestSaveAs_AutoFlush_Error(t *testing.T) {
	origFunc := execCommandFunc
	defer func() { execCommandFunc = origFunc }()

	execCommandFunc = func(name string, args ...string) *exec.Cmd {
		return exec.CommandContext(context.Background(), "false")
	}

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
	if err == nil {
		t.Fatal("Save() with failing flush should return error")
	}

	var fe *FlushError
	if !errors.As(err, &fe) {
		t.Errorf("Expected *FlushError from Save(), got %T: %v", err, err)
	}
}
