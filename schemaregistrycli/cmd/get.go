package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// get can handle three argument styles: <id>, <subj ver> or <subj>
var getCmd = &cobra.Command{
	Use:   "get <id> | (<subject> [<version>])",
	Short: "retrieves a schema specified by id or subject",
	Long: `The schema can be requested by id or subject.
When a subject is given, optionally one can provide a specific version. If no
version is specified, the latest version is returned.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("expected 1 to 2 arguments")
		}
		id, idParseErr := strconv.Atoi(args[0])
		var err error
		switch {
		case len(args) == 1 && idParseErr == nil:
			err = getById(id)
		case len(args) == 1 && idParseErr != nil:
			err = getLatestBySubject(args[0])
		case len(args) == 2:
			ver, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("2nd argument must be a version number")
			}
			err = getBySubjectVersion(args[0], ver)
		default:
			return fmt.Errorf("?")
		}
		return err
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
