package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var existsCmd = &cobra.Command{
	Use:   "exists <subject>",
	Short: "checks if the schema provided through stdin exists for the subject",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument")
		}
		isreg, sch, err := assertClient().IsRegistered(args[0], stdinToString())
		if err != nil {
			return err
		}
		fmt.Printf("exists: %v\n", isreg)
		if isreg {
			fmt.Printf("id: %d\n", sch.Id)
			fmt.Printf("version: %d\n", sch.Version)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(existsCmd)
}
