// Package word provides utilities for word games.
package main

import (
	"flag"
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/interpreter"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func debug(format string, tokens ...interface{}) {
	format = format + "\n"
	fmt.Fprintf(os.Stderr, format, tokens...)
}

func parseSpec(filename string) (interface{}, error) {
	f, _ := os.Open(filename)
	return dsl.ParseReader(filename, f, dsl.GlobalStore("filename", filename), dsl.Recover(false))
}

func fileDoesNotExist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func main() {
	outputFile := flag.String("dest", "entities.json", "destination file for generated content")
	flag.Parse()
	if len(os.Args) < 2 {
		log.Fatal("You must pass in a file")
	}

	filename := os.Args[len(os.Args)-1]
	if fileDoesNotExist(filename) {
		log.Fatalf("File passed '%v' does not exist\n", filename)
	}

	if tree, err := parseSpec(filename); err != nil {
		log.Fatalf("Error parsing %s: %v", filename, err)
	} else {
		if errors := interpreter.New(*outputFile).Visit(tree.(dsl.Node)); errors != nil {
			log.Fatalln(errors)
		}
	}
}
