package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:          "add <subject>",
	Short:        "registers the schema provided through stdin",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument")
		}
		id, err := assertClient(registryUrl).RegisterNewSchema(args[0], stdinToString())
		if err != nil {
			return err
		}
		log.Printf("registered schema with id %d\n", id)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
