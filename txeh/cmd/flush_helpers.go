package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/txn2/txeh"
)

// saveHosts handles the DryRun/Save/Flush pattern used by all mutation commands.
// On flush failure it prints a warning to stderr and does not exit with an error.
func saveHosts() {
	if DryRun {
		fmt.Print(etcHosts.RenderHostsFile())
		return
	}

	err := etcHosts.Save()
	if err != nil {
		var flushErr *txeh.FlushError
		if errors.As(err, &flushErr) {
			fmt.Fprintf(os.Stderr, "Warning: hosts file saved but DNS cache flush failed: %s\n", flushErr)
			if !Quiet {
				fmt.Fprintln(os.Stderr, "DNS cache may be stale. You can flush manually.")
			}
			return
		}
		fmt.Fprintf(os.Stderr, "Error: could not save %s. Reason: %s\n", etcHosts.WriteFilePath, err)
		os.Exit(1)
	}

	if Flush && !Quiet {
		fmt.Println("DNS cache flushed.")
	}
}
