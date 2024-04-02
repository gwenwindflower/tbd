package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDbtProfile(t *testing.T) {
	// Create a temporary profiles.yml file for testing
	tmpDir := t.TempDir()
	err := os.Mkdir(filepath.Join(tmpDir, ".dbt"), 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary .dbt directory: %v", err)
	}
	profilesFile := filepath.Join(tmpDir, ".dbt", "profiles.yml")
	content := []byte(`
test_profile:
  target: dev
  outputs:
    dev:
      type: snowflake
      account: testaccount
      user: testuser
      password: testpassword
      database: testdb
      warehouse: testwh
      schema: testschema
`)
	err = os.WriteFile(profilesFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary profiles.yml file: %v", err)
	}

	os.Setenv("HOME", tmpDir)

	// Profile exists
	profile, err := GetDbtProfile("test_profile")
	if err != nil {
		t.Errorf("GetDbtProfile returned an error for an existing profile: %v", err)
	}
	if profile.Target != "dev" {
		t.Errorf("Expected target 'dev', got '%s'", profile.Target)
	}

	// Profile does not exist
	profile, err = GetDbtProfile("aragorn")
	if err == nil {
		t.Error("GetDbtProfile did not return an error for a non-existing profile")
	}
	if profile != nil {
		t.Error("GetDbtProfile returned a non-nil profile for a non-existing profile")
	}
}
