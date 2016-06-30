package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var subjectsCmd = &cobra.Command{
	Use:   "subjects",
	Short: "lists all registered subjects",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		subs, err := assertClient().Subjects()
		if err != nil {
			return err
		}
		log.Printf("there are %d subjects\n", len(subs))
		for _, s := range subs {
			fmt.Println(s)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(subjectsCmd)
}
