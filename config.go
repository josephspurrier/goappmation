package goappmation

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// PortableInfo contains the settings for the portable application
type PortableInfo struct {
	ApplicationName   string            `json:"ApplicationName"`
	DownloadExtension string            `json:"DownloadExtension"`
	Version           string            `json:"Version"`
	VersionCheck      VersionCheck      `json:"VersionCheck"`
	RemoveRootFolder  bool              `json:"RemoveRootFolder"`
	RootFolderName    string            `json:"RootFolderName"`
	DownloadUrl       string            `json:"DownloadUrl"`
	ExtractRegExList  []string          `json:"ExtractRegExList"`
	CreateFolders     []string          `json:"CreateFolders"`
	CreateFiles       map[string]string `json:"CreateFiles"`
	MoveObjects       map[string]string `json:"MoveObjects"`
}

type VersionCheck struct {
	Url              string `json:"Url"`
	RegEx            string `json:"RegEx"`
	UseLatestVersion bool   `json:"UseLatestVersion"`
}

// ParseJSON parses the given bytes
func (pi *PortableInfo) parseJSON(jsonBytes []byte) error {
	return json.Unmarshal([]byte(jsonBytes), &pi)
}

// LoadConfig returns the struct from the config file
func LoadConfig(configFile string) (*PortableInfo, error) {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		return nil, err
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		return nil, err
	}

	// Create a new container
	pi := &PortableInfo{}

	// Parse the config
	if err := pi.parseJSON(jsonBytes); err != nil {
		return nil, err
	}

	return pi, nil
}
