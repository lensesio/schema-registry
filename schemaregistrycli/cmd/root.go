package cmd

import (
	"fmt"
	"os"

	"github.com/rollulus/schemaregistry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var registryUrl string = schemaregistry.DefaultUrl

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "schemaregistrycli",
	Short: "A command line interface for the Confluent schema registry",
	Long:  `A command line interface for the Confluent schema registry`,
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
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&registryUrl, "url", "e", schemaregistry.DefaultUrl, "schema registry url")

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.schemaregistrycli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".schemaregistrycli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")              // adding home directory as first search path
	viper.AutomaticEnv()                      // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
