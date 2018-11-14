package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// compatible can handle two argument styles: <subj ver> or <subj>
var compatibleCmd = &cobra.Command{
	Use:   "compatible <subject> [version]",
	Short: "tests compatibility between a schema from stdin and a given subject",
	Long: `The compatibility level of the subject is used for this check.
If it has never been changed, the global compatibility level applies.
If no schema version is specified, the latest version is tested.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("expected 1 to 2 arguments")
		}
		var iscompat bool
		var err error
		switch len(args) {
		case 1:
			iscompat, err = assertClient().IsLatestSchemaCompatible(args[0], stdinToString())
		case 2:
			ver, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("2nd argument must be a version number")
			}
			iscompat, err = assertClient().IsSchemaCompatible(args[0], stdinToString(), ver)
		}
		if err != nil {
			return err
		}
		if iscompat {
			fmt.Println("the provided schema is compatible")
		} else {
			err = fmt.Errorf("the provided schema is not compatible")
		}
		return err
	},
}

func init() {
	RootCmd.AddCommand(compatibleCmd)
}
