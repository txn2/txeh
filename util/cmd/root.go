package cmd

import (
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/txn2/txeh"
)

var rootCmd = &cobra.Command{
	Use:   "txeh",
	Short: "txeh is a /etc/hosts manager",
	Long: ` _            _
| |___  _____| |__
| __\ \/ / _ \ '_ \
| |_ >  <  __/ | | |
 \__/_/\_\___|_| |_| v` + Version + `

Add, remove and re-associate hostname entries in your /etc/hosts file.
Read more including useage as a Go library at https://github.com/txn2/txeh`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Printf("Error: can not display help, reason: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Please specify a sub-command such as \"add\" or \"remove\"")
		os.Exit(1)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEtcHosts()
	},
}

var (
	// HostsFileReadPath specify host file to read
	HostsFileReadPath string
	// HostsFileWritePath specify path to write resulting host file
	HostsFileWritePath string
	// Quiet results in no output
	Quiet bool
	// DryRun sends output to STDOUT (ignores quiet)
	DryRun bool

	etcHosts      *txeh.Hosts
	hostnameRegex *regexp.Regexp
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dryrun", "d", false, "dry run, output to stdout (ignores quiet)")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "no output")
	rootCmd.PersistentFlags().StringVarP(&HostsFileReadPath, "read", "r", "", "(override) Path to read /etc/hosts file.")
	rootCmd.PersistentFlags().StringVarP(&HostsFileWritePath, "write", "w", "", "(override) Path to write /etc/hosts file.")

	// validate hostnames (allow underscore for service records)
	hostnameRegex = regexp.MustCompile(`^([A-Za-z]|[0-9]|-|_|\.)+$`)
}

func validateCIDRs(cidrs []string) (bool, string) {
	for _, cr := range cidrs {
		if validateCIDR(cr) == false {
			return false, cr
		}
	}

	return true, ""
}

func validateCIDR(c string) bool {
	_, _, err := net.ParseCIDR(c)
	if err != nil {
		return false
	}

	return true
}

func validateIPAddresses(ips []string) (bool, string) {
	for _, ip := range ips {
		if validateIPAddress(ip) == false {
			return false, ip
		}
	}

	return true, ""
}

func validateIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	}

	return true
}

func validateHostnames(hostnames []string) (bool, string) {
	for _, hn := range hostnames {
		if validateHostname(hn) != true {
			return false, hn
		}
	}

	return true, ""
}

func validateHostname(hostname string) bool {
	return hostnameRegex.MatchString(hostname)
}

func emptyFilePaths() bool {
	return HostsFileReadPath == "" && HostsFileWritePath == ""
}

func initEtcHosts() {
	var (
		hosts *txeh.Hosts
		err   error
	)

	if emptyFilePaths() {
		hosts, err = txeh.NewHostsDefault()
	} else {
		hosts, err = txeh.NewHosts(&txeh.HostsConfig{
			ReadFilePath:  HostsFileReadPath,
			WriteFilePath: HostsFileWritePath,
		})
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	etcHosts = hosts
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}
