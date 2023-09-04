package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [IP] [HOSTNAME] [HOSTNAME] [HOSTNAME]...",
	Short: "Add hostnames to /etc/hosts",
	Long:  `Add/Associate one or more hostnames to an IP address `,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("the \"add\" command requires an IP address and at least one hostname")
		}

		if validateIPAddress(args[0]) == false {
			return errors.New("the IP address provided is not a valid ipv4 or ipv6 address")
		}

		if ok, hn := validateHostnames(args[1:]); !ok {
			return errors.New(fmt.Sprintf("\"%s\" is not a valid hostname", hn))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Adding host(s) \"%s\" to IP address %s\n", strings.Join(args[1:], " "), args[0])
		}

		AddHosts(args[0], args[1:])
	},
}

func AddHosts(ip string, hosts []string) {

	etcHosts.AddHosts(ip, hosts)

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
