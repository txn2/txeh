package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var exactHost bool

func init() {
	listCmd.AddCommand(listByHostsCmd)
	listByHostsCmd.PersistentFlags().BoolVarP(&exactHost, "exact", "e", false, "use -e for exact match")
}

var listByHostsCmd = &cobra.Command{
	Use:   "host [hostname] [hostname] [hostname]...",
	Short: "List IP for one or more hostnames",
	Long:  `List IP for one or more hostnames from /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"list host\" command requires at least one hostname")
		}

		if ok, cidr := validateHostnames(args); !ok {
			return errors.New(fmt.Sprintf("\"%s\" is not a valid hostname", cidr))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ListByHostnames(args)
	},
}

func ListByHostnames(hostnames []string) {
	for _, hn := range hostnames {
		addrHosts := etcHosts.ListAddressesByHost(hn, exactHost)
		for _, addrHost := range addrHosts {
			fmt.Printf("%s %s\n", addrHost[0], addrHost[1])
		}
	}
}
