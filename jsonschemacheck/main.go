package main

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"os"
)

func main() {

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + os.Args[1])
	documentLoader := gojsonschema.NewReferenceLoader("file://" + os.Args[2])

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("OK\n")
	} else {
		fmt.Printf("errors:\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}
