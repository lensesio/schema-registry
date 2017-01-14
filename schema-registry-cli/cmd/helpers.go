package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/datamountaineer/schema-registry"
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
	log.Printf("id: %d\n", sch.Id)
	var indented bytes.Buffer
	if err := json.Indent(&indented, []byte(sch.Schema), "", "  "); err != nil {
		fmt.Println(sch.Schema) //isn't a json object, which is legal
		return
	}
	indented.WriteTo(os.Stdout)
	os.Stdout.WriteString("\n")
}

func getById(id int) error {
	cl := assertClient()
	sch, err := cl.GetSchemaById(id)
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

func assertClient() schemaregistry.Client {
	c, err := schemaregistry.NewClient(viper.GetString("url"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return c
}
