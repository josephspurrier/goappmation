// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

package main

import (
	"github.com/josephspurrier/goappmation"
)

func main() {

	configFile := "../../config/MySQL Portable v5.7.9.json"

	// Run the automation
	goappmation.Run(configFile)
}
