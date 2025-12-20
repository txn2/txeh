package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listByCommentCmd)
}

var listByCommentCmd = &cobra.Command{
	Use:   "bycomment [COMMENT]",
	Short: "List hosts for a comment",
	Long:  `List all hostnames that have the specified comment in /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"list bycomment\" command requires a comment")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ListByComment(args[0])
	},
}

func ListByComment(comment string) {
	hosts := etcHosts.ListHostsByComment(comment)
	for _, h := range hosts {
		fmt.Println(h)
	}
}
