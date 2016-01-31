// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

package main

import (
	"flag"
	"log"
	"os"

	"github.com/josephspurrier/goappmation"
)

func main() {

	// Overwrite version
	flagVersion := flag.String("version", "", "Overwrites the version in the config file")

	flag.Parse()

	configFile := flag.Arg(0)

	if configFile == "" {
		log.Println("JSON Config file must be passed")
		os.Exit(1)
	}

	log.Println("TEST", *flagVersion)

	// Run the automation
	goappmation.Run(configFile, *flagVersion)
}
