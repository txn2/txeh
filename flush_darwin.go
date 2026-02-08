//go:build darwin

package txeh

import "runtime"

// flushDNSCachePlatform flushes the DNS cache on macOS.
// It runs dscacheutil -flushcache, then attempts killall -HUP mDNSResponder.
// The killall step is non-fatal since mDNSResponder may not be running.
func flushDNSCachePlatform() error {
	cmd := execCommandFunc("dscacheutil", "-flushcache") // #nosec G204 -- hardcoded command, not user input
	if err := cmd.Run(); err != nil {
		return &FlushError{
			Platform: runtime.GOOS,
			Command:  joinArgs("dscacheutil", []string{"-flushcache"}),
			Err:      err,
		}
	}

	// Attempt to restart mDNSResponder. This may fail if the process
	// is not running, which is not an error condition.
	cmd = execCommandFunc("killall", "-HUP", "mDNSResponder") // #nosec G204 -- hardcoded command, not user input
	_ = cmd.Run()

	return nil
}
