package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rollulus/schemaregistry"
)

func stdinToString() string {
	bs, err := ioutil.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func getById(id int) error {
	cl := assertClient(registryUrl)
	sch, err := cl.GetSchemaById(id)
	if err != nil {
		return err
	}
	fmt.Println(sch.Schema)
	return nil
}

func getLatestBySubject(subj string) error {
	cl := assertClient(registryUrl)
	sch, err := cl.GetLatestSchema(subj)
	if err != nil {
		return err
	}
	log.Printf("version: %d\n", sch.Version)
	log.Printf("id: %d\n", sch.Id)
	fmt.Println(sch.Schema.Schema)
	return nil
}

func getBySubjectVersion(subj string, ver int) error {
	cl := assertClient(registryUrl)
	sch, err := cl.GetSchemaBySubjectVersion(subj, ver)
	if err != nil {
		return err
	}
	log.Printf("version: %d\n", sch.Version)
	log.Printf("id: %d\n", sch.Id)
	fmt.Println(sch.Schema.Schema)
	return nil
}

func assertClient(endpoint string) *schemaregistry.Client {
	c, err := schemaregistry.NewClient(registryUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return c
}
