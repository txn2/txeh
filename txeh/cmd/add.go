package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addComment string

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addComment, "comment", "c", "", "Add an inline comment (e.g., 'managed-by-myapp')")
}

var addCmd = &cobra.Command{
	Use:   "add [IP] [HOSTNAME] [HOSTNAME] [HOSTNAME]...",
	Short: "Add hostnames to /etc/hosts",
	Long: `Add/Associate one or more hostnames to an IP address.

Use the --comment flag to add an inline comment that will appear after
the hostnames on the line (e.g., "127.0.0.1 myhost # my-comment").

Examples:
  txeh add 127.0.0.1 myhost
  txeh add 127.0.0.1 myhost --comment "managed-by-myapp"
  txeh add 127.0.0.1 svc1 svc2 svc3 -c "kubefwd"`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("the \"add\" command requires an IP address and at least one hostname")
		}

		if !validateIPAddress(args[0]) {
			return errors.New("the IP address provided is not a valid ipv4 or ipv6 address")
		}

		if ok, hn := validateHostnames(args[1:]); !ok {
			return fmt.Errorf("\"%s\" is not a valid hostname", hn)
		}

		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		if !Quiet {
			if addComment != "" {
				fmt.Printf("Adding host(s) \"%s\" to IP address %s with comment \"%s\"\n", strings.Join(args[1:], " "), args[0], addComment)
			} else {
				fmt.Printf("Adding host(s) \"%s\" to IP address %s\n", strings.Join(args[1:], " "), args[0])
			}
		}

		AddHosts(args[0], args[1:], addComment)
	},
}

// AddHosts adds hostnames to an IP address with an optional comment.
func AddHosts(ip string, hosts []string, comment string) {
	if comment != "" {
		etcHosts.AddHostsWithComment(ip, hosts, comment)
	} else {
		etcHosts.AddHosts(ip, hosts)
	}

	saveHosts()
}
