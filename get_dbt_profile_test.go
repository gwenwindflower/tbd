package main

import (
	"os"
	"testing"
)

func TestGetDbtProfile(t *testing.T) {
	CreateTempDbtProfile(t)
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	
  // Profile exists
	profile, err := GetDbtProfile("elf")
	if err != nil {
		t.Errorf("GetDbtProfile returned an error for an existing profile: %v", err)
	}
	if profile.Target != "dev" {
		t.Errorf("Expected target 'dev', got '%s'", profile.Target)
	}
	if profile.Outputs["dev"].ConnType != "snowflake" {
		t.Errorf("Expected connection type 'snowflake', got '%s'", profile.Outputs["dev"].ConnType)
	}

  // Profile exists, DuckDB
	profile, err = GetDbtProfile("dwarf")
	if err != nil {
		t.Errorf("GetDbtProfile returned an error for an existing profile: %v", err)
	}
	if profile.Target != "dev" {
		t.Errorf("Expected target 'dev', got '%s'", profile.Target)
	}
	if profile.Outputs["dev"].ConnType != "duckdb" {
		t.Errorf("Expected connection type 'duckdb', got '%s'", profile.Outputs["dev"].ConnType)
	}
	if profile.Outputs["dev"].Schema != "balins_tomb" {
		t.Errorf("Expected schema 'balins_tomb', got '%s'", profile.Outputs["dev"].Schema)
	}
  // If using dbt profile with DuckDB, path should be unedited
	if profile.Outputs["dev"].Path != "/usr/local/var/dwarf.db" {
		t.Errorf("Expected path '/usr/local/var/dwarf.db', got '%s'", profile.Outputs["dev"].Path)
	}

	// Profile does not exist
	profile, err = GetDbtProfile("dunedain")
	if err == nil {
		t.Error("GetDbtProfile did not return an error for a non-existing profile")
	}
	if profile != nil {
		t.Error("GetDbtProfile returned a non-nil profile for a non-existing profile")
	}
}
