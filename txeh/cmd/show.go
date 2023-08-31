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
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ShowHosts()
	},
}

func ShowHosts() {
	fmt.Print(etcHosts.RenderHostsFile())
}
