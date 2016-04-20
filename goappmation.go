package goappmation

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/josephspurrier/goappmation/bytesize"
)

func init() {
	// Verbose logging with file name and line number
	//log.SetFlags(log.Lshortfile)
}

var ForceExtract = true
var SkipDownload = true

// Run will execute commands from a config file
func Run(configFile string, versionOverwrite string) {

	// Set the name of the config file
	configFileName := path.Base(configFile)

	// Output the name of the config file
	log.Println("*** Loading: " + configFileName)

	// Load the config file
	pi, err := LoadConfig(configFile)
	if err != nil {
		log.Println(err)
		unifiedExit(1)
	}

	// Compile the regular expressions into one
	re, err := combineRegex(pi.ExtractRegExList)
	if err != nil {
		log.Println("Error creating regular expression from list |", err)
		unifiedExit(1)
	}

	// Set the application version
	version := pi.Version

	// Overwrite the version if available
	if versionOverwrite != "" {
		version = versionOverwrite
	}

	// If Version Check parameters are specified
	if pi.VersionCheck.Url != "" && pi.VersionCheck.RegEx != "" {
		// Extract the version from the webpage
		newVersion, err := extractFromRequest(pi.VersionCheck.Url, pi.VersionCheck.RegEx)
		if err != nil {
			log.Println("Error retrieving page |", err)
			unifiedExit(1)
		}

		if newVersion != version {
			log.Println("Newer version available: " + newVersion)
		}

		if pi.VersionCheck.UseLatestVersion {
			version = newVersion
		}

		log.Println("Using version: " + version)
	}

	// Set the folder name
	var folderName = strings.Replace(pi.ApplicationName, "{{VERSION}}", version, -1)

	// Set the zip name based off the folder
	// Note: The original file download name will be changed
	var zip = folderName + pi.DownloadExtension

	// If the zip file DOES exist on disk
	if isExist(zip) {
		// Output the filename of the folder
		log.Println("Download Exists:", zip)
	}

	// If SkipDownload is true
	if SkipDownload && isExist(zip) {
		log.Println("Skipping download")
	} else {
		log.Println("Will Download:", folderName)

		downloadURL := strings.Replace(pi.DownloadUrl, "{{VERSION}}", version, -1)

		log.Println("Downloading from:", downloadURL)
		log.Println("Downloading to:", zip)

		size, err := downloadFile(downloadURL, zip)
		if err != nil {
			log.Println("Error download file |", err)
			unifiedExit(1)
		}
		log.Println("Download Size:", bytesize.ByteSize(size))
	}

	// If the folder exists
	if isExist(folderName) {
		if ForceExtract {
			log.Println("Removing old folder:", folderName)
			err = os.RemoveAll(folderName)
			if err != nil {
				log.Println("Error removing working folder |", err)
				unifiedExit(1)
			}
		} else {
			log.Println("Folder already exists:", folderName)
			log.Println("*** No change")
		}
	}

	// Working folder is the root folder where the files will be extracted
	workingFolder := folderName

	// Root folder is directory relative to the current directory where the files
	// will be extracted to
	rootFolder := ""

	log.Println("Extracting files")

	switch pi.DownloadExtension {
	case ".zip":
		// If RemoveRootFolder is set to true
		if pi.RemoveRootFolder {
			// If the root folder name is specified
			if len(pi.RootFolderName) > 0 {
				workingFolder = pi.RootFolderName
			} else { // Else the root folder name is not specified so guess it
				// Return the name of the root folder in the ZIP
				workingFolder, err = extractZipRootFolder(zip)
				if err != nil {
					log.Println("Error discovering working folder |", err)
					unifiedExit(1)
				}
			}
		} else {
			rootFolder = workingFolder
		}

		// Extract files based on regular expression
		_, err = extractZipRegex(zip, rootFolder, re)
		if err != nil {
			log.Println("Error extracting from zip |", err)
			unifiedExit(1)
		}
	case ".msi":
		// Make the folder
		err = os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			log.Println("Error making folder |", err)
			unifiedExit(1)
		}

		// Get the full folder path
		fullFolderPath, err := filepath.Abs(folderName)
		if err != nil {
			log.Println("Error getting folder full path |", err)
			unifiedExit(1)
		}

		// Build the command
		cmd := exec.Command("msiexec")

		// Manually set the arguments since Go escaping does not work with MSI arguments
		argString := fmt.Sprintf(`/a "%v" /qb TARGETDIR="%v"`, zip, fullFolderPath)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    false,
			CmdLine:       " " + argString,
			CreationFlags: 0,
		}

		err = cmd.Run()
		if err != nil {
			log.Println("Error extracting from msi |", err)
			unifiedExit(1)
		}

		// If RemoveRootFolder is set to true
		if pi.RemoveRootFolder {
			// If the root folder name is specified
			if len(pi.RootFolderName) > 0 {

				//Get the full path of the folder to set as the root folder
				currentPath := filepath.Join(fullFolderPath, pi.RootFolderName)

				// Check to make sure the path is valid
				if currentPath == fullFolderPath {
					log.Println("RootFolderName is invalid:", pi.RootFolderName)
					unifiedExit(1)
				}

				// Copy files based on regular expressions
				_, err = copyMsiRegex(currentPath, fullFolderPath+"_temp", re)
				if err != nil {
					log.Println("Error copy from msi folder |", err)
					unifiedExit(1)
				}

				// Set the working folder so the rename will work later
				workingFolder = fullFolderPath + "_temp"

				// Remove the original full folder path
				err = os.RemoveAll(fullFolderPath)
				if err != nil {
					log.Println("Error removing MSI folder:", currentPath)
					unifiedExit(1)
				}

			} else { // Else the root folder name is not specified
				log.Println("The string, RemoveRootName, is required for MSIs")
				unifiedExit(1)
			}
		} else {
			log.Println("The boolean, RemoveRootFolder, is required for MSIs")
			unifiedExit(1)
		}
	default:
		log.Println("Download extension not supported:", pi.DownloadExtension)
		unifiedExit(1)
	}

	log.Println("Creating folders")
	err = createFolders(pi.CreateFolders, workingFolder)
	if err != nil {
		log.Println("Error creating folders |", err)
		unifiedExit(1)
	}

	log.Println("Creating files")
	err = writeScripts(pi.CreateFiles, workingFolder)
	if err != nil {
		log.Println("Error writing files |", err)
		unifiedExit(1)
	}

	log.Println("Moving objects")
	err = moveObjects(pi.MoveObjects, workingFolder)
	if err != nil {
		log.Println("Error moving objects |", err)
		unifiedExit(1)
	}

	log.Println("Renaming folder to:", folderName)
	err = os.Rename(workingFolder, folderName)
	if err != nil {
		log.Println("Error renaming folder |", err)
		unifiedExit(1)
	}

	//unifiedExit(0)
	log.Println("*** Success")
}
