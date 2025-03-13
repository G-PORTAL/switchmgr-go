package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadTestData reads the testdata file for the given command. Because we can't
// use the Driver interface in the tests, we have to mock the output of the
// commands. This function reads the testdata file for the given command and
// returns the content as string.
func ReadTestData(command string, jsonMessage interface{}) error {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(directory + "/testdata")
	for os.IsNotExist(err) {
		directory, _ = filepath.Abs(directory + "/..")
		_, err = os.Stat(directory + "/testdata/arista_eos/")
	}

	// add version to mockfile name
	mockfile := strings.Replace(command, " ", "-", -1)

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s.json", directory+"/testdata/arista_eos/", mockfile))
	if err != nil {
		return err
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, jsonMessage)
}
