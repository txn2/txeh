// Package cmd implements the txeh command-line interface.
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
 \__/_/\_\___|_| |_| Version: ` + VersionFromBuild() + `

Add, remove and re-associate hostname entries in your /etc/hosts file.
Read more including usage as a Go library at https://github.com/txn2/txeh`,
	Run: func(cmd *cobra.Command, _ []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Printf("Error: can not display help, reason: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Please specify a sub-command such as \"add\" or \"remove\"")
		os.Exit(1)
	},
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		initEtcHosts()
	},
}

var (
	// HostsFileReadPath specifies the host file to read.
	HostsFileReadPath string
	// HostsFileWritePath specifies the path to write the resulting host file.
	HostsFileWritePath string
	// Quiet results in no output.
	Quiet bool
	// DryRun sends output to STDOUT (ignores quiet).
	DryRun bool
	// Flush triggers a DNS cache flush after writing the hosts file.
	Flush bool
	// MaxHostsPerLine limits hostnames per line (0=auto, -1=unlimited, >0=explicit).
	MaxHostsPerLine int

	etcHosts      *txeh.Hosts
	hostnameRegex *regexp.Regexp
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dryrun", "d", false, "dry run, output to stdout (ignores quiet)")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "no output")
	rootCmd.PersistentFlags().StringVarP(&HostsFileReadPath, "read", "r", "", "(override) Path to read /etc/hosts file.")
	rootCmd.PersistentFlags().StringVarP(&HostsFileWritePath, "write", "w", "", "(override) Path to write /etc/hosts file.")
	rootCmd.PersistentFlags().BoolVarP(&Flush, "flush", "f", false, "flush DNS cache after modifying hosts file")
	rootCmd.PersistentFlags().IntVarP(&MaxHostsPerLine, "max-hosts-per-line", "m", 0, "Max hostnames per line (0=auto, -1=unlimited, >0=explicit). Auto uses 9 on Windows.")

	// validate hostnames (allow underscore for service records)
	// disallow leading dots, trailing dots, and consecutive dots
	hostnameRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+(\.[A-Za-z0-9_-]+)*$`)
}

func validateCIDRs(cidrs []string) (ok bool, failed string) {
	for _, cr := range cidrs {
		if !validateCIDR(cr) {
			return false, cr
		}
	}

	return true, ""
}

func validateCIDR(c string) bool {
	_, _, err := net.ParseCIDR(c)
	return err == nil
}

func validateIPAddresses(ips []string) (ok bool, failed string) {
	for _, ip := range ips {
		if !validateIPAddress(ip) {
			return false, ip
		}
	}

	return true, ""
}

func validateIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func validateHostnames(hostnames []string) (ok bool, failed string) {
	for _, hn := range hostnames {
		if !validateHostname(hn) {
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
	if os.Getenv("TXEH_AUTO_FLUSH") == "1" {
		Flush = true
	}

	var (
		hosts *txeh.Hosts
		err   error
	)

	if emptyFilePaths() && MaxHostsPerLine == 0 && !Flush {
		hosts, err = txeh.NewHostsDefault()
	} else {
		hosts, err = txeh.NewHosts(&txeh.HostsConfig{
			ReadFilePath:    HostsFileReadPath,
			WriteFilePath:   HostsFileWritePath,
			MaxHostsPerLine: MaxHostsPerLine,
			AutoFlush:       Flush,
		})
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	etcHosts = hosts
}

// Execute bootstraps the cobra root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
