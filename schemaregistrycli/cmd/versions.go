package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "lists all available versions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("expected 1 argument")
		}
		client := assertClient(registryUrl)
		vers, err := client.Versions(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", vers)
	},
}

func init() {
	RootCmd.AddCommand(versionsCmd)
}
