package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	Version = "0.3"
)

var version = flag.Bool("version", false, "show gofsp version")
var filename = flag.String("file", "", "policy file")

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("gofsp - %s\n", Version)
		os.Exit(0)
	}

	fsp := NewFspServer()

	// load policy from file
	if len(*filename) > 0 {
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatalln(err)
		}
		fsp.LoadPolicy(file)
	}

	fsp.Startup()
}
