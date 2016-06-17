package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <subject>",
	Short: "registers the schema provided through stdin",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("expected 1 argument")
		}
		id, err := assertClient(registryUrl).RegisterNewSchema(args[0], stdinToString())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("registered schema with id %d\n", id)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
