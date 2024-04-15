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

func TestSetConnectionDetailsWithDuckDBDbtProfile(t *testing.T) {
	CreateTempDbtProfile(t)
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	formResponse := FormResponse{
		UseDbtProfile:        true,
		DbtProfile:           "dwarf",
		DbtProfileOutput:     "dev",
		Schema:               "balins_tomb",
		GenerateDescriptions: false,
		BuildDir:             "test_build",
		Confirm:              true,
	}
	connectionDetails := SetConnectionDetails(formResponse)
	want := shared.ConnectionDetails{
		ConnType: "duckdb",
		Path:     "/usr/local/var/dwarf.db",
		Database: "khazad_dum",
		Schema:   "balins_tomb",
	}
	if connectionDetails != want {
		t.Errorf("got %v, want %v", connectionDetails, want)
	}
}

func TestSetConnectionDetailsWithDuckDBWithoutDbtProfile(t *testing.T) {
	formResponse := FormResponse{
		UseDbtProfile:        false,
		Warehouse:            "duckdb",
		Path:                 "dwarf.db",
		Database:             "khazad_dum",
		Schema:               "balins_tomb",
		GenerateDescriptions: false,
		BuildDir:             "test_build",
		Confirm:              true,
	}
	connectionDetails := SetConnectionDetails(formResponse)
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get working directory: %v", err)
	}
	p := wd + "/dwarf.db"
	want := shared.ConnectionDetails{
		ConnType: "duckdb",
		Path:     p,
		Database: "khazad_dum",
		Schema:   "balins_tomb",
	}
	if connectionDetails != want {
		t.Errorf("got %v, want %v", connectionDetails, want)
	}
}
