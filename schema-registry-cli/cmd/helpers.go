package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hokaccha/go-prettyjson"

	schemaregistry "github.com/landoop/schema-registry"
	"github.com/spf13/viper"
)

func stdinToString() string {
	bs, err := ioutil.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func printSchema(sch schemaregistry.Schema) {
	log.Printf("version: %d\n", sch.Version)
	log.Printf("id: %d\n", sch.ID)

	pretty, err := prettyjson.Format([]byte(sch.Schema))
	if err != nil {
		fmt.Println(sch.Schema) //isn't a json object, which is legal
		return
	}
	os.Stdout.Write(pretty)
	os.Stdout.WriteString("\n")
}

func getByID(id int) error {
	cl := assertClient()
	sch, err := cl.GetSchemaByID(id)
	if err != nil {
		return err
	}
	fmt.Println(sch)
	return nil
}

func getLatestBySubject(subj string) error {
	cl := assertClient()
	sch, err := cl.GetLatestSchema(subj)
	if err != nil {
		return err
	}
	printSchema(sch)
	return nil
}

func getBySubjectVersion(subj string, ver int) error {
	cl := assertClient()
	sch, err := cl.GetSchemaBySubject(subj, ver)
	if err != nil {
		return err
	}
	printSchema(sch)
	return nil
}

func printConfig(cfg schemaregistry.Config, subj string) {
	if subj == "" {
		subj = "global"
	}
	if cfg.CompatibilityLevel == "" {
		cfg.CompatibilityLevel = "not defined, using global"
	}
	fmt.Printf("%s compatibility-level: %s\n", subj, cfg.CompatibilityLevel)
}

func getConfig(subj string) error {
	cl := assertClient()
	cfg, err := cl.GetConfig(subj)
	if err != nil {
		return err
	}
	printConfig(cfg, subj)
	return nil
}

func assertClient() *schemaregistry.Client {
	c, err := schemaregistry.NewClient(viper.GetString("url"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	username := viper.GetString("username")
	password := viper.GetString("password")

	if username != "" && password != "" {
		err = c.SetBasicAuth(username, password)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	return c
}
