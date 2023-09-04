package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var Version = "0.0.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of txeh",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		version := VersionFromBuild()
		fmt.Printf("txeh Version: %s\n", version)
	},
}

// Version returns the version of txeh binary
func VersionFromBuild() (version string) {
	// Version is managed with goreleaser
	if Version != "0.0.0" {
		return Version
	}
	// Version is managed by "go install"
	b, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if b == nil {
		version = "nil"
	} else {
		version = b.Main.Version
	}
	return version
}
