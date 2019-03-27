package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
	schemaregistry "github.com/landoop/schema-registry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	registryURL string
	verbose     bool
	nocolor     bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "schema-registry-cli",
	Short: "A command line interface for the Confluent schema registry",
	Long:  `A command line interface for the Confluent schema registry`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !verbose {
			log.SetOutput(ioutil.Discard)
		}
		if nocolor {
			color.NoColor = true
		}
		log.Printf("schema registry url: %s\n", viper.Get("url"))
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "be verbose")
	RootCmd.PersistentFlags().BoolVarP(&nocolor, "no-color", "n", false, "dont color output")
	RootCmd.PersistentFlags().StringVarP(&registryURL, "url", "e", schemaregistry.DefaultURL, "schema registry url, overrides SCHEMA_REGISTRY_URL")
	viper.SetEnvPrefix("schema_registry")
	viper.BindPFlag("url", RootCmd.PersistentFlags().Lookup("url"))
	viper.BindEnv("url")
}
