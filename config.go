package lassie

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	// DefaultAddr is the default address of the Congress API. You normally won't
	// have to change this.
	DefaultAddr = "https://api.lora.telenor.io"

	// ConfigFile is the default name for the config file. The configuration
	// file is a plain text file that contains the default Lassie configuration.
	// The configuration file is expected to be at the current home directory.
	ConfigFile = ".lassie"

	// AddressEnvironmentVariable is the name of the environment variable that
	// can be used to override the address set in the configuration file.
	// If the  environment variable isn't set (or is empty) the configuration
	// file settings will be used.
	AddressEnvironmentVariable = "LASSIE_ADDRESS"

	// TokenEnvironmentVariable is the name of the environment variable that
	// can be used to override the token set in the configuration file.
	TokenEnvironmentVariable = "LASSIE_TOKEN"
)

// These are the configuration file directives.

const (
	addressKey = "address"
	tokenKey   = "token"
)

// Return both address and token from configuration file. The file name is
// for testing purposes; use the ConfigFile constant when calling the functino.
func addressTokenFromConfig(filename string) (string, string) {
	address, token := readConfig(getFullPath(filename))

	envAddress := os.Getenv(AddressEnvironmentVariable)
	if envAddress != "" {
		address = envAddress
	}

	envToken := os.Getenv(TokenEnvironmentVariable)
	if envToken != "" {
		token = envToken
	}

	return address, token
}

func getFullPath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(usr.HomeDir, filename)
}

// readConfig reads the config file and returns the address and token
// settings from the file.
func readConfig(filename string) (string, string) {
	address := DefaultAddr
	token := ""

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return address, token
	}
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(bufio.ScanLines)
	lineno := 0
	for scanner.Scan() {
		lineno++
		line := strings.ToLower(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			// ignore comments and empty lines
			continue
		}
		words := strings.Split(scanner.Text(), "=")
		if len(words) != 2 {
			fmt.Printf("Not a key value expression on line %d in %s: %s\n", lineno, filename, scanner.Text())
			continue
		}
		switch words[0] {
		case addressKey:
			address = strings.TrimSpace(words[1])
		case tokenKey:
			token = strings.TrimSpace(words[1])
		default:
			fmt.Printf("Unknown keyword on line %d in %s: %s\n", lineno, filename, scanner.Text())
		}
	}
	return address, token
}
