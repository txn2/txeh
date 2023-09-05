package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.AddCommand(removeHostCmd)
}

var removeHostCmd = &cobra.Command{
	Use:   "host [HOSTNAME] [HOSTNAME] [HOSTNAME]...",
	Short: "Remove hostnames from /etc/hosts",
	Long:  `Remove one or more hostnames from /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"remove hosts\" command requires at least one hostname to remove")
		}

		if ok, hn := validateHostnames(args); !ok {
			return errors.New(fmt.Sprintf("\"%s\" is not a valid hostname", hn))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Removing host(s) \"%s\"\n", strings.Join(args, " "))
		}

		RemoveHosts(args)
	},
}

func RemoveHosts(hosts []string) {

	etcHosts.RemoveHosts(hosts)

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
