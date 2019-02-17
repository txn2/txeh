package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.AddCommand(removeCidrCmd)
}

var removeCidrCmd = &cobra.Command{
	Use:   "cidr [CIDR] [CIDR] [CIDR]...",
	Short: "Remove ranges of addresses from /etc/hosts",
	Long:  `Remove one or more CIDR ranges of IP addresses from /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"remove cidr\" command requires at least one CIDR range to remove")
		}

		if ok, ip := validateCIDRs(args); !ok {
			return errors.New(fmt.Sprintf("\"%s\" is not a valid CIDR", ip))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Removing ip ranges(s) \"%s\"\n", strings.Join(args, " "))
		}

		RemoveIPRanges(args)
	},
}

func RemoveIPRanges(cidrs []string) {

	addresses := make([]string, 0)

	// loop through all the CIDR ranges
	for _, cidr := range cidrs {

		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			fmt.Printf("Error: there was a problem with the IP range %s. Reason: %s\n", cidr, err.Error())
			os.Exit(1)
		}

		hfLines := etcHosts.GetHostFileLines()

		for _, hfl := range *hfLines {

			ip := net.ParseIP(hfl.Address)
			if ip != nil {
				if ipnet.Contains(ip) {
					addresses = append(addresses, hfl.Address)
				}
			}
		}

	}

	etcHosts.RemoveAddresses(addresses)

	if DryRun {
		fmt.Print(etcHosts.RenderHostsFile())
		return
	}

	err := etcHosts.Save()
	if err != nil {
		fmt.Printf("Error: could not save %s. Reason: %s\n", etcHosts.WriteFilePath, err.Error())
		os.Exit(1)
	}
}
