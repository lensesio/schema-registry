package cmd

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("expected 1 argument")
		}
		bs, err := ioutil.ReadAll(bufio.NewReader(os.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s\n", string(bs))
		id, err := assertClient(registryUrl).RegisterNewSchema(args[0], string(bs))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("registered schema with id %d\n", id)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
