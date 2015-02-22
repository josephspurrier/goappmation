// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/josephspurrier/goappmation"
)

func main() {

	configFile := "../../config/mysql.json"

	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		log.Printf("Cannot open %q: %v", configFile, err)
		os.Exit(1)
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		log.Printf("Error reading %q: %v", configFile, err)
		os.Exit(1)
	}

	// Create a new container
	pi := &goappmation.PortableInfo{}

	// Parse the config
	if err := pi.ParseJSON(jsonBytes); err != nil {
		log.Printf("Could not parse the .json file: %v", err)
		os.Exit(2)
	}

	// Run the automation
	goappmation.Run(pi)
}
