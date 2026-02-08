package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.AddCommand(removeCommentCmd)
}

var removeCommentCmd = &cobra.Command{
	Use:   "bycomment [COMMENT]",
	Short: "Remove all hosts with a comment",
	Long:  `Remove all host entries that have the specified comment from /etc/hosts`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"remove bycomment\" command requires a comment")
		}

		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Removing all hosts with comment \"%s\"\n", args[0])
		}

		RemoveByComment(args[0])
	},
}

// RemoveByComment removes all host entries with the given comment.
func RemoveByComment(comment string) {
	etcHosts.RemoveByComment(comment)
	saveHosts()
}
