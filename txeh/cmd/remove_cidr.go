package cmd

import (
	"errors"
	"fmt"
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

		if ok, cidr := validateCIDRs(args); !ok {
			return fmt.Errorf("\"%s\" is not a valid CIDR", cidr)
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
	err := etcHosts.RemoveCIDRs(cidrs)
	if err != nil {
		fmt.Printf("Error: there was a problem parsing a CIDR. Reason: %s\n", err.Error())
		os.Exit(1)
	}

	if DryRun {
		fmt.Print(etcHosts.RenderHostsFile())
		return
	}

	err = etcHosts.Save()
	if err != nil {
		fmt.Printf("Error: could not save %s. Reason: %s\n", etcHosts.WriteFilePath, err.Error())
		os.Exit(1)
	}
}
