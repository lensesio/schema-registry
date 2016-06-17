package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var existsCmd = &cobra.Command{
	Use:   "exists <subject>",
	Short: "checks if the schema provided through stdin exists for the subject",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("expected 1 argument")
		}
		isreg, sch, err := assertClient(registryUrl).IsRegistered(args[0], stdinToString())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("exists: %v\n", isreg)
		if !isreg {
			return
		}
		fmt.Printf("id: %d\n", sch.Id)
		fmt.Printf("version: %d\n", sch.Version)
	},
}

func init() {
	RootCmd.AddCommand(existsCmd)
}
