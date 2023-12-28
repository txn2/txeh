package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update [Original IP] [New IP], [HOSTNAME] [HOSTNAME] [HOSTNAME]...",
	Short: "Update hostnames in /etc/hosts",
	Long:  `Update/Associate one or more hostnames to an IP address `,
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) < 3 {
			return errors.New("the \"update\" command requires an Original IP address, a New IP address and at least one hostname")
		}

		if validateIPAddress(args[0]) == false {
			return errors.New("the Original IP address provided is not a valid ipv4 or ipv6 address")
		}

		if validateIPAddress(args[1]) == false {
			return errors.New("the New IP address provided is not a valid ipv4 or ipv6 address")
		}

		var hosts []string
		for _, host := range args[2:] {

			for i := 0; i < len(host); i++ {
				if host[i] == ':' {
					break
				}
				hosts = append(hosts, host)
			}
		}

		if ok, hn := validateHostnames(hosts); !ok {
			return errors.New(fmt.Sprintf("\"%s\" is not a valid hostname", hn))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Updating host(s) \"%s\" from IP address %s to %s\n", strings.Join(args[2:], " "), args[0], args[1])
		}

		var hosts []string
		var comments []string

		for _, arg := range args[2:] {
			// If we see a colon, we're done with hosts and we're on to comments joined by spaces
			if arg[0] == ':' {
				// Split everything before the colon and remove it.
				comments = append(comments, strings.Split(strings.Join(args[2:], " "), ":")[1])
				break
			}
			hosts = append(hosts, arg)
		}

		UpdateHosts(args[0], args[1], hosts, comments...)

	},
}

func UpdateHosts(oldIP string, newIP string, hosts []string, comment ...string) {

	etcHosts.UpdateHosts(oldIP, newIP, hosts, comment...)

	if DryRun {
		fmt.Println("-----HOSTS-----")
		fmt.Println(hosts)
		fmt.Println("-----HOSTS-----")
		fmt.Println("-----COMMENTS-----")
		fmt.Println(comment)
		fmt.Println("-----COMMENTS-----")
		fmt.Println("-----ETC HOSTS-----")
		fmt.Print(etcHosts.RenderHostsFile())
		fmt.Println("-----ETC HOSTS-----")
		return
	}

	err := etcHosts.Save()
	if err != nil {
		fmt.Printf("Error: could not save %s. Reason: %s\n", etcHosts.WriteFilePath, err.Error())
		os.Exit(1)
	}
}
