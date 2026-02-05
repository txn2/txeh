package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list [BY_TYPE] [IP|CIDR] [IP|CIDR] [IP|CIDR] ...",
	Short: "List hostnames or IP addresses",
	Long:  `List hostnames or IP addresses present in /etc/hosts`,
	Run: func(cmd *cobra.Command, _ []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Printf("Error: can not display help, reason: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Please specify a sub-command such as 'ip', 'cidr' or 'host'")
		os.Exit(1)
	},
}
