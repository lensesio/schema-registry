package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var subjectsCmd = &cobra.Command{
	Use:   "subjects",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		subs, err := assertClient(registryUrl).Subjects()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", subs)
	},
}

func init() {
	RootCmd.AddCommand(subjectsCmd)
}
