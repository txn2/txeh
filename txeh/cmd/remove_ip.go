package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.AddCommand(removeIPCmd)
}

var removeIPCmd = &cobra.Command{
	Use:   "ip [IP] [IP] [IP]...",
	Short: "Remove IP addresses from /etc/hosts",
	Long:  `Remove one or more IP addresses from /etc/hosts`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"remove ip\" command requires at least one IP address to remove")
		}

		if ok, ip := validateIPAddresses(args); !ok {
			return fmt.Errorf("\"%s\" is not a valid ip address", ip)
		}

		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Removing ip(s) \"%s\"\n", strings.Join(args, " "))
		}

		removeIPs(args)
	},
}

// removeIPs removes the given IP addresses from the hosts file.
func removeIPs(ips []string) {
	etcHosts.RemoveAddresses(ips)
	saveHosts()
}
