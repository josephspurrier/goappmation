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

var msiDstFolder string
var msiSrcFolder string
var msiAllowRegExp *regexp.Regexp

// copyMsiRegex will copy certain files from a directory to another folder based on a regular expression
func copyMsiRegex(srcFolder string, dstFolder string, allowRegExp *regexp.Regexp) (bool, error) {
	// Create folder to extract files
	if !isExist(dstFolder) {
		err := os.MkdirAll(dstFolder, os.ModePerm)
		if err != nil {
			return false, err
		}
	}

	msiDstFolder = dstFolder
	msiSrcFolder = srcFolder
	msiAllowRegExp = allowRegExp

	err := filepath.Walk(srcFolder, visitMsiFile)
	if err != nil {
		return false, err
	}

	return true, nil
}

func visitMsiFile(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		return nil // can't walk here, but continue walking elsewhere
	}

	// Path AFTER the source directory (not including the src dir)
	relativePath := strings.TrimLeft(fp, msiSrcFolder)

	// Destination path
	finalPath := filepath.Join(msiDstFolder, relativePath)

	// Destination path folder
	basePath := filepath.Dir(finalPath)

	// Check if the file matches the regular expression
	if !msiAllowRegExp.MatchString(strings.Replace(relativePath, "\\", "/", -1)) {
		return nil
	}

	// Create the file directory if it doesn't exist
	if !isExist(basePath) {
		err = os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Move the file
	err = os.Rename(fp, finalPath)
	if err != nil {
		return err
	}

	return nil
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
			err = os.MkdirAll(rootFolder, os.ModePerm)
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
