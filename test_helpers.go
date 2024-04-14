package main

import (
	"os"
	"path/filepath"
	"testing"
)

func CreateTempDbtProfile(t *testing.T) string {
	content := []byte(`
elf:
  target: dev
  outputs:
    dev:
      type: snowflake
      account: 123456.us-east-1
      user: legolas
      database: mirkwood
      schema: hall_of_thranduil
      threads: 8

human:
  target: dev
  outputs:
    dev:
      type: bigquery
      method: oauth
      project: gondor
      dataset: minas_tirith
      threads: 16

dwarf:
  target: dev
  outputs:
    dev:
      type: duckdb
      path: /usr/local/var/dwarf.db
      database: khazad_dum
      schema: balins_tomb
      threads: 4
`)
	tmpDir := t.TempDir()
	err := os.Mkdir(filepath.Join(tmpDir, ".dbt"), 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary .dbt directory: %v", err)
	}
	profilesFile := filepath.Join(tmpDir, ".dbt", "profiles.yml")
	err = os.WriteFile(profilesFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary profiles.yml file: %v", err)
	}
	os.Setenv("HOME", tmpDir)
	return tmpDir
}
