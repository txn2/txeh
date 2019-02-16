package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.AddCommand(removeIpCmd)
}

var removeIpCmd = &cobra.Command{
	Use:   "ip [IP] [IP] [IP]...",
	Short: "Remove IP addresses from /etc/hosts",
	Long:  `Remove one or more IP addresses from /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"remove ip\" command requires at least one IP address to remove")
		}

		if ok, ip := validateIPAddresses(args); !ok {
			return errors.New(fmt.Sprintf("\"%s\" is not a valid ip address", ip))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Removing ip(s) \"%s\"\n", strings.Join(args, " "))
		}

		RemoveIPs(args)
	},
}

func RemoveIPs(ips []string) {

	etcHosts.RemoveAddresses(ips)

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
