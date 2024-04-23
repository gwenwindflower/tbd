package internal

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
	ps, err := FetchDbtProfiles()
	if err != nil {
		t.Errorf("Error fetching dbt profiles: %v", err)
	}
	connectionDetails := SetConnectionDetails(formResponse, ps)
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
	CreateTempDbtProfiles(t)
	ps, err := FetchDbtProfiles()
	if err != nil {
		t.Errorf("Error fetching dbt profiles: %v", err)
	}
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	formResponse := FormResponse{
		UseDbtProfile:        true,
		DbtProfileName:       "elf",
		DbtProfileOutput:     "dev",
		Database:             "mirkwood",
		Schema:               "hall_of_thranduil",
		GenerateDescriptions: false,
		BuildDir:             "test_build",
		Confirm:              true,
	}
	connectionDetails := SetConnectionDetails(formResponse, ps)
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
	CreateTempDbtProfiles(t)
	ps, err := FetchDbtProfiles()
	if err != nil {
		t.Errorf("Error fetching dbt profiles: %v", err)
	}
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	formResponse := FormResponse{
		UseDbtProfile:        true,
		DbtProfileName:       "dwarf",
		DbtProfileOutput:     "dev",
		Database:             "khazad_dum",
		Schema:               "balins_tomb",
		GenerateDescriptions: false,
		BuildDir:             "test_build",
		Confirm:              true,
	}
	connectionDetails := SetConnectionDetails(formResponse, ps)
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
	ps, err := FetchDbtProfiles()
	if err != nil {
		t.Errorf("Error fetching dbt profiles: %v", err)
	}
	connectionDetails := SetConnectionDetails(formResponse, ps)
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

func TestSetConnectionDetailsPostgresWithoutDbtProfile(t *testing.T) {
	formResponse := FormResponse{
		UseDbtProfile:        false,
		Warehouse:            "postgres",
		Host:                 "localhost",
		Port:                 "5432",
		Username:             "treebeard",
		Password:             "entmoot",
		Database:             "fangorn",
		Schema:               "huorns",
		SslMode:              "disable",
		GenerateDescriptions: false,
		BuildDir:             "test_build",
		Confirm:              true,
	}
	ps, err := FetchDbtProfiles()
	if err != nil {
		t.Errorf("Error fetching dbt profiles: %v", err)
	}
	cd := SetConnectionDetails(formResponse, ps)
	want := shared.ConnectionDetails{
		ConnType: "postgres",
		Host:     "localhost",
		Port:     5432,
		Username: "treebeard",
		Password: "entmoot",
		Database: "fangorn",
		Schema:   "huorns",
		SslMode:  "disable",
	}
	if cd != want {
		t.Errorf("got %v, want %v", cd, want)
	}
}
