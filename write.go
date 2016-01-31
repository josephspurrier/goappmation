package goappmation

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// writeScripts creates files from the config file
func writeScripts(scripts map[string]string, workingFolder string) error {
	// Loop through each script
	for name, body := range scripts {

		// Path of file
		relativePath := filepath.Join(workingFolder, name)

		// Write to file
		err := ioutil.WriteFile(relativePath, []byte(body), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// moveObjects moves or renames files
func moveObjects(files map[string]string, workingFolder string) error {
	// Loop through each file
	for dst, src := range files {

		// Path of files
		dstFile := filepath.Join(workingFolder, dst)
		srcFile := filepath.Join(workingFolder, src)

		// Rename Files
		err := os.Rename(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}
