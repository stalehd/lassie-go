package lassie

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileDefaultConfig(t *testing.T) {

	contents := "address=http://example.com\ntoken=sometoken"
	tempFile := "lassie.testconfig"
	ioutil.WriteFile(getFullPath(tempFile), []byte(contents), 0666)

	// unset the environment first to make sure it won't interfere with the
	// file
	oldAddr := os.Getenv(AddressEnvironmentVariable)
	oldToken := os.Getenv(TokenEnvironmentVariable)
	defer func() {
		os.Setenv(AddressEnvironmentVariable, oldAddr)
		os.Setenv(TokenEnvironmentVariable, oldToken)
	}()

	os.Setenv(AddressEnvironmentVariable, "")
	os.Setenv(TokenEnvironmentVariable, "")

	address, token := addressTokenFromConfig(tempFile)
	if address != "http://example.com" || token != "sometoken" {
		t.Fatalf("Configuration isn't the expected values: %s / %s", address, token)
	}

	contents = "token=foobar\nsome=thing\nother=thing\n\n\n"
	ioutil.WriteFile(getFullPath(tempFile), []byte(contents), 0666)
	address, token = addressTokenFromConfig(tempFile)
	if address != DefaultAddr || token != "foobar" {
		t.Fatalf("Configuration isn't the expected values: %s / %s", address, token)
	}

	os.Remove(getFullPath(tempFile))
}

func TestEnvironmentConfig(t *testing.T) {
	oldAddr := os.Getenv(AddressEnvironmentVariable)
	oldToken := os.Getenv(TokenEnvironmentVariable)
	defer func() {
		os.Setenv(AddressEnvironmentVariable, oldAddr)
		os.Setenv(TokenEnvironmentVariable, oldToken)
	}()

	os.Setenv(AddressEnvironmentVariable, "something")
	os.Setenv(TokenEnvironmentVariable, "other")

	address, token := addressTokenFromConfig(ConfigFile)

	if address != "something" {
		t.Fatal("Expected environment variable to override config")
	}
	if token != "other" {
		t.Fatal("Expected environment variable to override config")
	}
}
