package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadTestData reads the testdata file for the given command. Because we can't
// use the Driver interface in the tests, we have to mock the output of the
// commands. This function reads the testdata file for the given command and
// returns the content as string.
func ReadTestData(command string, version *int) string {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(directory + "/testdata")
	for os.IsNotExist(err) {
		directory, _ = filepath.Abs(directory + "/..")
		_, err = os.Stat(directory + "/testdata/fsos/")
	}

	// add version to mockfile name
	mockfile := strings.Replace(command, " ", "-", -1)
	if version != nil {
		mockfile = fmt.Sprintf("%s-%v", mockfile, *version)
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s.txt", directory+"/testdata/fsos/", mockfile))
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return ""
	}
	return string(data)
}
