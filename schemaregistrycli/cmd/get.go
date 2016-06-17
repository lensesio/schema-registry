package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// get can handle three argument styles: <id>, <subj ver> or <subj>
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 2 {
			log.Fatalf("expected 1 to 2 arguments")
		}
		id, idParseErr := strconv.Atoi(args[0])
		switch {
		case len(args) == 1 && idParseErr == nil:
			getById(id)
		case len(args) == 1 && idParseErr != nil:
			getLatestBySubject(args[0])
		case len(args) == 2:
			ver, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("2nd argument must be a version number")
			}
			getBySubjectVersion(args[0], ver)
		default:
			log.Fatalf("?")
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
