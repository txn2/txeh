package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of txeh",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("txeh Version %s", VERSION)
	},
}
