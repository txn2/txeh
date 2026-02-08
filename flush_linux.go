//go:build linux

package txeh

import (
	"errors"
	"os/exec"
	"runtime"
)

// ErrNoResolver is returned when no supported DNS resolver is found on the system.
var ErrNoResolver = errors.New("no supported DNS resolver found (tried resolvectl, systemd-resolve)")

// flushDNSCachePlatform flushes the DNS cache on Linux.
// It tries resolvectl first, then falls back to systemd-resolve.
// Returns ErrNoResolver (wrapped in FlushError) if neither is available.
func flushDNSCachePlatform() error {
	// Try resolvectl first (systemd 239+).
	if _, err := exec.LookPath("resolvectl"); err == nil {
		cmd := execCommandFunc("resolvectl", "flush-caches") // #nosec G204 -- hardcoded command, not user input
		if err := cmd.Run(); err != nil {
			return &FlushError{
				Platform: runtime.GOOS,
				Command:  joinArgs("resolvectl", []string{"flush-caches"}),
				Err:      err,
			}
		}
		return nil
	}

	// Fallback to systemd-resolve (older systemd).
	if _, err := exec.LookPath("systemd-resolve"); err == nil {
		cmd := execCommandFunc("systemd-resolve", "--flush-caches") // #nosec G204 -- hardcoded command, not user input
		if err := cmd.Run(); err != nil {
			return &FlushError{
				Platform: runtime.GOOS,
				Command:  joinArgs("systemd-resolve", []string{"--flush-caches"}),
				Err:      err,
			}
		}
		return nil
	}

	return &FlushError{
		Platform: runtime.GOOS,
		Command:  "",
		Err:      ErrNoResolver,
	}
}
