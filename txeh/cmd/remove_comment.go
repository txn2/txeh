package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.AddCommand(removeCommentCmd)
}

var removeCommentCmd = &cobra.Command{
	Use:   "bycomment [COMMENT]",
	Short: "Remove all hosts with a comment",
	Long:  `Remove all host entries that have the specified comment from /etc/hosts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("the \"remove bycomment\" command requires a comment")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !Quiet {
			fmt.Printf("Removing all hosts with comment \"%s\"\n", args[0])
		}

		RemoveByComment(args[0])
	},
}

func RemoveByComment(comment string) {
	etcHosts.RemoveByComment(comment)

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
