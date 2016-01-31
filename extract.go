package goappmation

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// extractFromRequest will return extracted text from a page at a URL
func extractFromRequest(url string, regExp string) (string, error) {
	body, err := getRequestBody(url)
	if err != nil {
		return "", err
	}

	re, err := regexp.Compile(regExp)
	if err != nil {
		return "", err
	}

	rePair := re.FindStringSubmatch(body)

	if len(rePair) > 1 {
		version := strings.TrimSpace(rePair[1])

		return version, nil
	}

	return "", errors.New("Could not find string on page:" + url)
}

// extractZipRegex will extract certain files from a ZIP file to a folder based on a regular expression
func extractZipRegex(file string, rootFolder string, allowRegExp *regexp.Regexp) (bool, error) {
	// Open a zip archive
	r, err := zip.OpenReader(file)
	if err != nil {
		return false, err
	}
	defer r.Close()

	// If the rootFolder is NOT empty,
	if rootFolder != "" {
		// Create folder to extract files
		if !isExist(rootFolder) {
			os.MkdirAll(rootFolder, os.ModePerm)
			if err != nil {
				return false, err
			}
		}
	}

	// Loop through all files
	for _, f := range r.File {

		if !allowRegExp.MatchString(f.Name) {
			continue
		}

		// Path of file
		relativePath := filepath.Join(rootFolder, f.Name)

		// Path of file directory
		basePath := filepath.Dir(relativePath)

		// If the object is a directory, create it
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(relativePath, os.ModePerm)
			if err != nil {
				return false, err
			}
			continue
		}

		// Create the file directory if it doesn't exist
		if !isExist(basePath) {
			err = os.MkdirAll(basePath, os.ModePerm)
			if err != nil {
				return false, err
			}
		}

		rc, err := f.Open()
		if err != nil {
			return false, err
		}

		// Create the file
		out, err := os.Create(relativePath)
		defer out.Close()
		if err != nil {
			return false, err
		}

		// Write the file
		_, err = io.Copy(out, rc)
		if err != nil {
			return false, err
		}
		rc.Close()
	}

	return true, nil
}

// extractZipRootFolder will extract a folder from a ZIP file
func extractZipRootFolder(file string) (string, error) {
	// Open a zip archive
	r, err := zip.OpenReader(file)
	if err != nil {
		return "", err
	}
	defer r.Close()

	if len(r.File) > 0 {
		pathArray := strings.Split(r.File[0].Name, "/")
		return pathArray[0], nil
	}

	return "", errors.New("Working folder not found in first file path")
}

// Unused **********************************************************************

// extractZipAll will extract all files from a ZIP file to a folder
func extractZipAll(file string, workingFolder string) (bool, error) {
	// Open a zip archive
	r, err := zip.OpenReader(file)
	if err != nil {
		return false, err
	}
	defer r.Close()

	// Create folder to extract files
	if !isExist(workingFolder) {
		os.MkdirAll(workingFolder, os.ModePerm)
		if err != nil {
			return false, err
		}
	}

	// Loop through all files
	for _, f := range r.File {
		//log.Printf("Extracting %s\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			return false, err
		}

		// Path of file
		relativePath := filepath.Join(workingFolder, f.Name)

		// Path of file directory
		basePath := filepath.Dir(relativePath)

		// Create the file directory if it doesn't exist
		if !isExist(basePath) {
			err = os.MkdirAll(basePath, os.ModePerm)
			if err != nil {
				return false, err
			}
		}

		// Create the file
		out, err := os.Create(relativePath)
		defer out.Close()
		if err != nil {
			return false, err
		}

		// Write the file
		_, err = io.Copy(out, rc)
		if err != nil {
			return false, err
		}
		rc.Close()
	}

	return true, nil
}
