package main

import (
	"os"
	"path"
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
      authenticator: externalbrowser
      account: dunedain.snowflakecomputing.com
      user: aragorn
      database: gondor
      schema: minas_tirith
      threads: 8
`)
	tpp := path.Join(tmpDir, "profiles.yml")
	got, err := os.ReadFile(tpp)
	if err != nil {
		t.Fatalf("Failed to read profiles.yml: %v", err)
	}
	// os.Remove("profiles.yml")
	if string(got) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}
