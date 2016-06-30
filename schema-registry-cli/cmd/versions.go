package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "lists all available versions",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument")
		}
		client := assertClient()
		vers, err := client.Versions(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", vers)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionsCmd)
}
