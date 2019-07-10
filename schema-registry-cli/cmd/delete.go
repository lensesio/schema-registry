package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

// get can handle three argument styles: <id>, <subj ver> or <subj>
var deleteCmd = &cobra.Command{
	Use:   "delete <subject>",
	Short: "delete a subject",
	Long:``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument")
		}

		subject := args[0]

		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Warning! You are deleting the subject '%s'. Are you sure", subject),
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			fmt.Printf("Aborted.\n")
			return nil
		}

		err = deleteSubject(subject)
		if err != nil {
			return err
		}
		fmt.Printf("'%s' has been deleted\n", subject)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
