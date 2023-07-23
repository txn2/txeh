package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listByCidrCmd)
}

var listByCidrCmd = &cobra.Command{
	Use:   "cidr [CIDR] [CIDR] [CIDR]...",
	Short: "List hosts for one or more CIDRs",
	Long:  `List hosts for one or more CIDRs from /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"list cidr\" command requires at least one CIDR address")
		}

		if ok, cidr := validateCIDRs(args); !ok {
			return errors.New(fmt.Sprintf("\"%s\" is not a valid CIDR", cidr))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ListByCIDRs(args)
	},
}

func ListByCIDRs(cidrs []string) {
	for _, cidr := range cidrs {
		ipHosts := etcHosts.ListHostsByCIDR(cidr)
		for _, ih := range ipHosts {
			fmt.Printf("%s %s %s\n", cidr, ih[0], ih[1])
		}
	}
}
