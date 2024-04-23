package internal

import (
	"os"
	"testing"
)

func TestGetDbtProfile(t *testing.T) {
	CreateTempDbtProfiles(t)
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	ps, err := FetchDbtProfiles()
	if err != nil {
		t.Errorf("Error fetching dbt profiles: %v", err)
	}
	// Profile exists
	p, err := GetDbtProfile("elf", ps)
	if err != nil {
		t.Errorf("GetDbtProfile returned an error for an existing profile: %v", err)
	}
	if p.Target != "dev" {
		t.Errorf("Expected target 'dev', got '%s'", p.Target)
	}
	if p.Outputs["dev"].ConnType != "snowflake" {
		t.Errorf("Expected connection type 'snowflake', got '%s'", p.Outputs["dev"].ConnType)
	}

	// Profile exists, DuckDB
	p, err = GetDbtProfile("dwarf", ps)
	if err != nil {
		t.Errorf("GetDbtProfile returned an error for an existing profile: %v", err)
	}
	if p.Target != "dev" {
		t.Errorf("Expected target 'dev', got '%s'", p.Target)
	}
	if p.Outputs["dev"].ConnType != "duckdb" {
		t.Errorf("Expected connection type 'duckdb', got '%s'", p.Outputs["dev"].ConnType)
	}
	if p.Outputs["dev"].Schema != "balins_tomb" {
		t.Errorf("Expected schema 'balins_tomb', got '%s'", p.Outputs["dev"].Schema)
	}
	// If using dbt profile with DuckDB, path should be unedited
	if p.Outputs["dev"].Path != "/usr/local/var/dwarf.db" {
		t.Errorf("Expected path '/usr/local/var/dwarf.db', got '%s'", p.Outputs["dev"].Path)
	}

	// Profile does not exist
	p, err = GetDbtProfile("dunedain", ps)
	if err == nil {
		t.Error("GetDbtProfile did not return an error for a non-existing profile")
	}
	if p.Outputs != nil {
		t.Error("GetDbtProfile returned a non-empty profile for a non-existing profile")
	}
}
