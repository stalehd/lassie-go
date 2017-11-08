package lassie

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ClientConfig is a client configuration for Lassie-go
type ClientConfig interface {
	// Address returns the address to use
	Address() string
	// Token returns the token to use
	Token() string
}

// NewWithConfig creates a new client with the provided configuration
func NewWithConfig(config ClientConfig) (*Client, error) {
	return NewWithAddr(config.Address(), config.Token())
}

// EnvironmentConfig gets its configuration from environment variables
type EnvironmentConfig struct {
}

// NewEnvironmentConfig creates a new EnvironmentConfig
func NewEnvironmentConfig() *EnvironmentConfig {
	return &EnvironmentConfig{}
}

// Address tries to retrieve the endpoint address (ie https://api.lora.telenor.io)
// from the environment variable LASSIE_ADDRESS. If is empty it will return the
// default.
func (e *EnvironmentConfig) Address() string {
	endpoint := os.Getenv("LASSIE_ADDRESS")
	if endpoint == "" {
		return DefaultAddr
	}
	return endpoint
}

// Token returns the token from the environment variable LASSIE_TOKEN
func (e *EnvironmentConfig) Token() string {
	return os.Getenv("LASSIE_TOKEN")
}

// UserConfig gets its configuration from a file named ~/.lassie
//
// The configuration format is quite simple with "key=value" entries on each
// line. Only two parameters are supported -- "Address" and "Token". The Address
// parameter is optional.
//
// Note: This isn't tested on Windows.
type UserConfig struct {
	values map[string]string
}

// NewUserConfig creates a new UserConfig
func NewUserConfig() *UserConfig {
	return &UserConfig{}
}

const (
	addressKey = "address"
	tokenKey   = "token"
)

func (u *UserConfig) configFile() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(usr.HomeDir, ".lassie")
}

// readConfig populates the values map with two elements -- "token" and
// "address"
func (u *UserConfig) readConfig(filename string) map[string]string {
	if u.values == nil {
		u.values = make(map[string]string)
		u.values[addressKey] = DefaultAddr
		u.values[tokenKey] = ""

		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			return u.values
		}
		scanner := bufio.NewScanner(bytes.NewReader(buf))
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			words := strings.Split(scanner.Text(), "=")
			if len(words) == 1 {
				continue
			}
			switch strings.ToLower(strings.TrimSpace(words[0])) {
			case addressKey:
				u.values[addressKey] = strings.TrimSpace(words[1])
			case tokenKey:
				u.values[tokenKey] = strings.TrimSpace(words[1])
			}
		}
	}
	return u.values
}

// Address reads the configuration file and checks if the file has the field
// "address". If the field isn't found it will use the default address. If the
// file is missing it will use the defaults
func (u *UserConfig) Address() string {
	u.readConfig(u.configFile())
	return u.values[addressKey]
}

// Token reads the API token from the field "token" in the configuration
// file. If the configuration file is missing it will return the defaults.
func (u *UserConfig) Token() string {
	u.readConfig(u.configFile())
	return u.values[tokenKey]
}