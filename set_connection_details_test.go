package main

import (
	"os"
	"testing"

	"github.com/gwenwindflower/tbd/shared"
)

func TestSetConnectionDetailsWithoutDbtProfile(t *testing.T) {
	formResponse := FormResponse{
		UseDbtProfile:        false,
		Warehouse:            "snowflake",
		Username:             "aragorn",
		Account:              "123456.us-east-1",
		Database:             "gondor",
		Schema:               "minas_tirith",
		GenerateDescriptions: false,
		BuildDir:             "test_build",
		Confirm:              true,
	}
	connectionDetails := SetConnectionDetails(formResponse)
	want := shared.ConnectionDetails{
		ConnType: "snowflake",
		Username: "aragorn",
		Account:  "123456.us-east-1",
		Database: "gondor",
		Schema:   "minas_tirith",
	}
	if connectionDetails != want {
		t.Errorf("got %v, want %v", connectionDetails, want)
	}
}

func TestSetConnectionDetailsWithDbtProfile(t *testing.T) {
	CreateTempDbtProfile(t)
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	formResponse := FormResponse{
		UseDbtProfile:        true,
		DbtProfile:           "elf",
		DbtProfileOutput:     "dev",
		Schema:               "hall_of_thranduil",
		GenerateDescriptions: false,
		BuildDir:             "test_build",
		Confirm:              true,
	}
	connectionDetails := SetConnectionDetails(formResponse)
	want := shared.ConnectionDetails{
		ConnType: "snowflake",
		Username: "legolas",
		Account:  "123456.us-east-1",
		Database: "mirkwood",
		Schema:   "hall_of_thranduil",
	}
	if connectionDetails != want {
		t.Errorf("got %v, want %v", connectionDetails, want)
	}
}
