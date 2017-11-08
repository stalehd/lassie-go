package lassie

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUserConfig(t *testing.T) {
	contents := "address=http://example.com\ntoken=sometoken"
	tempFile := "tempconfig"
	ioutil.WriteFile(tempFile, []byte(contents), 0666)

	uc := UserConfig{}
	uc.values = uc.readConfig(tempFile)
	if uc.Address() != "http://example.com" || uc.Token() != "sometoken" {
		t.Fatalf("Configuration isn't the expected values: %v", uc)
	}

	contents = "token=foobar\nsome=thing\nother=thing\n\n\n"
	ioutil.WriteFile(tempFile, []byte(contents), 0666)
	uc = UserConfig{}
	uc.values = uc.readConfig(tempFile)
	if uc.Address() != DefaultAddr || uc.Token() != "foobar" {
		t.Fatalf("Configuration isn't the expected values: %v", uc)
	}

	os.Remove(tempFile)
}

func TestEnvironmentConfig(t *testing.T) {
	os.Setenv("LASSIE_ADDRESS", "")
	os.Setenv("LASSIE_TOKEN", "")

	env := EnvironmentConfig{}
	if env.Address() != DefaultAddr {
		t.Fatal("Expected address to be default")
	}
	if env.Token() != "" {
		t.Fatal("Expected token to be empty")
	}

	os.Setenv("LASSIE_ADDRESS", "https://example.com")
	os.Setenv("LASSIE_TOKEN", "foo")
	if env.Address() != "https://example.com" {
		t.Fatal("Expected environment variable to override config")
	}
	if env.Token() != "foo" {
		t.Fatal("Expected environment variable to override config")
	}
}

func TestClientConfig(t *testing.T) {
	// Just exercise the configurations
	NewEnvironmentConfig()
	NewUserConfig()
	NewConfig()
}
