// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

package goappmation

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/josephspurrier/goappmation/bytesize"
)

var ForceExtract = false
var SkipDownload = false

type PortableInfo struct {
	ApplicationName     string
	DownloadExtension   string
	CreateWorkingFolder bool
	ExplicitVersion     string
	ExplicitFolderName  string
	VersionUrl          string
	VersionRegEx        string
	DownloadUrl         string
	ExtractRegExList    []string
	CreateFiles         map[string]string
}

// ParseJSON parses the given bytes
func (pi *PortableInfo) ParseJSON(jsonBytes []byte) error {
	return json.Unmarshal([]byte(jsonBytes), &pi)
}

func getRequestBody(url string) (string, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client, err := http.DefaultClient.Do(r)
	defer client.Body.Close()

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(client.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ExtractFromRequest(url string, regExp string) (string, error) {
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

func DownloadFile(url string, fileName string) (int64, error) {
	out, err := os.Create(fileName)
	defer out.Close()
	if err != nil {
		return 0, err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(out, resp.Body)

	return n, err
}

func ExtractZipAll(file string, workingFolder string) (bool, error) {
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

func ExtractZipRegex(file string, workingFolder string, allowRegExp *regexp.Regexp) (bool, error) {
	// Open a zip archive
	r, err := zip.OpenReader(file)
	if err != nil {
		return false, err
	}
	defer r.Close()

	if workingFolder != "" {
		// Create folder to extract files
		if !isExist(workingFolder) {
			os.MkdirAll(workingFolder, os.ModePerm)
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

func ExtractZipRootFolder(file string) (string, error) {
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

func isExist(dir string) bool {
	if _, err := os.Stat(dir); err == nil {
		return true
	}

	return false
}

func CombineRegex(s []string) (*regexp.Regexp, error) {
	joined := strings.Join(s, "|")

	re, err := regexp.Compile(joined)
	if err != nil {
		return nil, err
	}

	return re, nil
}

func WriteScripts(scripts map[string]string, workingFolder string) (bool, error) {
	// Loop through each script
	for name, body := range scripts {

		// Path of file
		relativePath := filepath.Join(workingFolder, name)

		// Write to file
		err := ioutil.WriteFile(relativePath, []byte(body), os.ModePerm)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func Run(pi *PortableInfo) {

	var version string
	var err error

	if pi.ExplicitVersion == "" {
		version, err = ExtractFromRequest(pi.VersionUrl, pi.VersionRegEx)
		if err != nil {
			log.Println("Error retrieving page |", err)
			os.Exit(1)
		}
	} else {
		version = pi.ExplicitVersion
	}

	var folderName = pi.ApplicationName + "v" + version
	var zip = folderName + pi.DownloadExtension

	if isExist(zip) {
		log.Println("Current:", folderName)
		if !ForceExtract {
			os.Exit(0)
		}
	} else {
		log.Println("New Version:", pi.ApplicationName+"v"+version)
	}

	if !SkipDownload {
		downloadURL := strings.Replace(pi.DownloadUrl, "{VERSION}", version, -1)

		log.Println("Downloading from:", downloadURL)
		log.Println("Downloading to:", zip)

		size, err := DownloadFile(downloadURL, zip)
		if err != nil {
			log.Println("Error download file |", err)
			os.Exit(1)
		}
		log.Println("Download Size:", bytesize.ByteSize(size))
	}

	re, err := CombineRegex(pi.ExtractRegExList)
	if err != nil {
		log.Println("Error creating regular express |", err)
		os.Exit(1)
	}

	var workingFolder = ""
	var passFolder = ""

	if pi.CreateWorkingFolder {
		workingFolder = pi.ApplicationName
		passFolder = workingFolder
	} else {
		workingFolder, err = ExtractZipRootFolder(zip)
		if err != nil {
			log.Println("Error discovering working folder |", err)
			os.Exit(1)
		}
		passFolder = ""
	}

	if isExist(folderName) {
		if ForceExtract {
			log.Println("Removing old folder:", folderName)
			err = os.RemoveAll(folderName)
			if err != nil {
				log.Println("Error removing working folder |", err)
				os.Exit(1)
			}
		} else {
			log.Println("Folder already exists:", folderName)
			os.Exit(1)
		}
	}

	log.Println("Extracting to:", workingFolder)

	_, err = ExtractZipRegex(zip, passFolder, re)
	if err != nil {
		log.Println("Error extracting from zip |", err)
		os.Exit(1)
	}

	log.Println("Creating files")
	WriteScripts(pi.CreateFiles, workingFolder)
	if err != nil {
		log.Println("Error writing files |", err)
		os.Exit(1)
	}

	// Get config folder name
	if pi.ExplicitFolderName != "" {
		folderName = pi.ExplicitFolderName
	}

	log.Println("Renaming folder to:", folderName)
	//os.RemoveAll(folderName)
	err = os.Rename(workingFolder, folderName)
	if err != nil {
		log.Println("Error renaming folder |", err)
		os.Exit(1)
	}

	log.Println("Complete")
}
