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
		log.Fatal(err)
	}
	return string(bs)
}

func getById(id int) {
	cl := assertClient(registryUrl)
	sch, err := cl.GetSchemaById(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sch.Schema)
}

func getLatestBySubject(subj string) {
	cl := assertClient(registryUrl)
	sch, err := cl.GetLatestSchema(subj)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("version: %d\n", sch.Version)
	log.Printf("id: %d\n", sch.Id)
	fmt.Println(sch.Schema.Schema)
}

func getBySubjectVersion(subj string, ver int) {
	cl := assertClient(registryUrl)
	sch, err := cl.GetSchemaBySubjectVersion(subj, ver)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("version: %d\n", sch.Version)
	log.Printf("id: %d\n", sch.Id)
	fmt.Println(sch.Schema.Schema)
}

func assertClient(endpoint string) *schemaregistry.Client {
	c, err := schemaregistry.NewClient(registryUrl)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
