package main

import (
	"flag"
	"fmt"
	"os"
)

const APP_NAME = "to"
const APP_VERSION = "0.0.1b"

func main() {
	optVersion := flag.Bool("version", false,
		"Show version.")
	flag.Parse()

	if *optVersion {
		fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 1 {
		// Too many args.
		os.Exit(1)
	}

	os.Exit(0)
}
