package goappmation

import (
	"log"
	"os"
	"regexp"
	"strings"
)

// isExist returns true if a file object exists
func isExist(dir string) bool {
	if _, err := os.Stat(dir); err == nil {
		return true
	}

	return false
}

// combineRegex will take a string array of regular expressions and compile them
// into a single regular expressions
func combineRegex(s []string) (*regexp.Regexp, error) {
	joined := strings.Join(s, "|")

	re, err := regexp.Compile(joined)
	if err != nil {
		return nil, err
	}

	return re, nil
}

// unifiedExit prints a line and then exists
func unifiedExit(exitCode int) {
	if exitCode == 0 {
		log.Println("*** Success")
		os.Exit(exitCode)
	} else {
		log.Println("*** Fail")
		os.Exit(exitCode)
	}
}
