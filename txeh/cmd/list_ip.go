package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listByIpCmd)
}

var listByIpCmd = &cobra.Command{
	Use:   "ip [IP] [IP] [IP]...",
	Short: "List hosts for one or more IP addresses",
	Long:  `List hosts for one or more IP addresses from /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"list ip\" command requires at least one IP address")
		}

		if ok, ip := validateIPAddresses(args); !ok {
			return fmt.Errorf("\"%s\" is not a valid ip address", ip)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ListByIPs(args)
	},
}

func ListByIPs(ips []string) {
	for _, ip := range ips {
		hosts := etcHosts.ListHostsByIP(ip)
		for _, h := range hosts {
			fmt.Printf("%s %s\n", ip, h)
		}
	}
}
