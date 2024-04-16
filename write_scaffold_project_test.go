package main

import (
	"os"
	"path"
	"testing"

	"github.com/gwenwindflower/tbd/shared"
)

func TestWriteScaffoldProject(t *testing.T) {
	cd := shared.ConnectionDetails{
		ConnType: "snowflake",
		Username: "user",
		Account:  "account",
		Database: "database",
		Schema:   "schema",
		Project:  "project",
	}
	bd := t.TempDir()
	pn := "project"
	_, err := WriteScaffoldProject(cd, bd, pn)
	if err != nil {
		t.Fatalf("Error scaffolding project: %v\n", err)
	}
	// Check that the directories were created
	for _, folder := range []string{"models", "analyses", "macros", "seeds", "snapshots", "data-tests", "models/staging", "models/marts"} {
		if _, err := os.Stat(path.Join(bd, folder)); os.IsNotExist(err) {
			t.Fatalf("Directory %s was not created\n", folder)
		}
	}
	// Check that .gitignore was created correctly
	gitignore := []byte(`.venv
venv
.env
env

target/
dbt_packages/
logs/
profiles.yml

.DS_Store

.user.yml

.ruff_cache
__pycache__
`)
	gi := path.Join(bd, ".gitignore")
	got, err := os.ReadFile(gi)
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}
	if string(got) != string(gitignore) {
		t.Errorf("Expected %s, got %s", gitignore, got)
	}
	// Check that project.yml was created correctly
	projectYaml := []byte(`config-version: 2

name: project
profile: snowflake

model-paths: ["models"]
analysis-paths: ["analyses"]
test-paths: ["data-tests"]
seed-paths: ["seeds"]
macro-paths: ["macros"]
snapshot-paths: ["snapshots"]

target-path: "target"
clean-targets:
  - "target"
  - "dbt_packages"

models:
  project:
    staging:
      +materialized: view
    marts:
      +materialized: table
`)
	py := path.Join(bd, "dbt_project.yml")
	got, err = os.ReadFile(py)
	if err != nil {
		t.Fatalf("Failed to read dbt_project.yml: %v", err)
	}
	if string(got) != string(projectYaml) {
		t.Errorf("Expected %s, got %s", projectYaml, got)
	}
}
