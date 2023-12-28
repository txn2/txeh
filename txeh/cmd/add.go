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
	Use:   "add [IP] [HOSTNAME] [HOSTNAME] [HOSTNAME]... :[COMMENT]...",
	Short: "Add hostnames to /etc/hosts",
	Long:  `Add/Associate one or more hostnames to an IP address `,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("the \"add\" command requires an IP address and at least one hostname")
		}

		if validateIPAddress(args[0]) == false {
			return errors.New("the IP address provided is not a valid ipv4 or ipv6 address")
		}

		var hosts []string
		for _, host := range args[1:] {

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
			fmt.Printf("Adding host(s) \"%s\" to IP address %s\n", strings.Join(args[1:], " "), args[0])
		}

		var hosts []string
		var comments []string

		for _, arg := range args[1:] {
			// If we see a colon, we're done with hosts and we're on to comments joined by spaces
			if arg[0] == ':' {
				comments = append(comments, strings.Split(strings.Join(args[2:], " "), ":")[1])
				break
			}
			hosts = append(hosts, arg)
		}

		AddHosts(args[0], hosts, comments...)
	},
}

func AddHosts(ip string, hosts []string, comment ...string) {

	etcHosts.AddHosts(ip, hosts, comment...)

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
