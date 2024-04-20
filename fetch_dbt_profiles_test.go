package main

import (
	"os"
	"testing"
)

func TestFetchDbtProfiles(t *testing.T) {
	CreateTempDbtProfiles(t)
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	profiles, err := FetchDbtProfiles()
	if err != nil {
		t.Fatalf("Error fetching dbt profiles: %v\n", err)
	}
	if err != nil {
		t.Fatalf("Error fetching dbt profiles: %v\n", err)
	}
	if profiles["elf"].Outputs["dev"].ConnType != "snowflake" {
		t.Fatalf("Expected snowflake, got %s\n", profiles["elf"].Outputs["dev"].ConnType)
	}
	if profiles["human"].Outputs["dev"].ConnType != "bigquery" {
		t.Fatalf("Expected bigquery, got %s\n", profiles["human"].Outputs["dev"].ConnType)
	}
	if profiles["dwarf"].Outputs["dev"].ConnType != "duckdb" {
		t.Fatalf("Expected duckdb, got %s\n", profiles["dwarf"].Outputs["dev"].ConnType)
	}
}
