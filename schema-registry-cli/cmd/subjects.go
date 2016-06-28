package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var subjectsCmd = &cobra.Command{
	Use:   "subjects",
	Short: "lists all registered subjects",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		subs, err := assertClient(registryUrl).Subjects()
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", subs)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(subjectsCmd)
}
