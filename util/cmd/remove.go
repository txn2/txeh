package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove [TYPE] [HOSTNAME|IP] [HOSTNAME|IP] [HOSTNAME|IP]...",
	Short: "Remove a hostname or ip address",
	Long:  `Remove one or more hostnames or ip addresses from /etc/hosts`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Printf("Error: can not display help, reason: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Please specify a sub-command such as \"host\" or \"ip\"")
		os.Exit(1)
	},
}
