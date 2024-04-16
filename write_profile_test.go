package main

import (
	"os"
	"testing"

	"github.com/gwenwindflower/tbd/shared"
)

func TestWriteProfile(t *testing.T) {
	cd := shared.ConnectionDetails{
		ConnType: "snowflake",
		Username: "aragorn",
		Account:  "dunedain.snowflakecomputing.com",
		Database: "gondor",
		Schema:   "minas_tirith",
	}
	tmpDir := t.TempDir()
	WriteProfile(cd, tmpDir)

	expected := []byte(`
snowflake:
  target: dev
  outputs:
    dev:
      type: snowflake
      account: dunedain.snowflakecomputing.com
      user: aragorn
      database: gondor
      schema: minas_tirith
      threads: 8
`)
	got, err := os.ReadFile("profiles.yml")
	if err != nil {
		t.Fatalf("Failed to read profiles.yml: %v", err)
	}
	// os.Remove("profiles.yml")
	if string(got) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}
