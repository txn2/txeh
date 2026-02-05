package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show hostnames in /etc/hosts",
	Long:  `Show the content of the /etc/hosts file `,
	Args: func(_ *cobra.Command, _ []string) error {
		return nil
	},
	Run: func(_ *cobra.Command, _ []string) {
		ShowHosts()
	},
}

// ShowHosts prints the current hosts file content.
func ShowHosts() {
	fmt.Print(etcHosts.RenderHostsFile())
}
