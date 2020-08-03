package cmd

import (
	"fmt"
	"log"

	schemaregistry "github.com/landoop/schema-registry"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:          "add <subject> [<schemaType>]",
	Short:        "registers the schema provided through stdin",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("expected 1 to 2 arguments")
		}
		var schemaType schemaregistry.SchemaType

		if len(args) == 1 {
			schemaType = schemaregistry.AVRO
		} else {
			switch args[1] {

			case "JSON":
				schemaType = schemaregistry.JSON
			case "PROTOBUF":
				schemaType = schemaregistry.PROTOBUF
			case "AVRO":
				schemaType = schemaregistry.AVRO
			default:
				return fmt.Errorf("schemaType must be one of AVRO, JSON, PROTOBUF")
			}
		}

		id, err := assertClient().RegisterNewSchemaV2(args[0], stdinToString(), schemaType)
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
