package main

import (
	"log"
	"os"
	"path"
	"text/template"

	"github.com/gwenwindflower/tbd/shared"
)

func WriteScaffoldProject(cd shared.ConnectionDetails, bd string, pn string) (string, error) {
	folders := []string{"models", "analyses", "macros", "seeds", "snapshots", "data-tests", "models/staging", "models/marts"}
	emptyFolders := []string{"analyses", "macros", "seeds", "snapshots", "data-tests", "models/marts"}
	for _, folder := range folders {
		p := path.Join(bd, folder)
		err := os.MkdirAll(p, 0755)
		if err != nil {
			return "", err
		}
	}
	for _, folder := range emptyFolders {
		p := path.Join(bd, folder, ".gitkeep")
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Fatalf("Failed to create .gitkeep in %s folder %v\n", folder, err)
		}
	}
	projectYamlTemplate := `config-version: 2

name: {{.ProjectName}}
profile: {{.ConnType}}

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
  {{.ProjectName}}:
    staging:
      +materialized: view
    marts:
      +materialized: table
`
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

	tmpl, err := template.New("dbt_project.yml").Parse(projectYamlTemplate)
	if err != nil {
		log.Fatalf("Failed to parse dbt_project.yml template %v\n", err)
	}
	p := path.Join(bd, "dbt_project.yml")
	o, err := os.Create(p)
	if err != nil {
		log.Fatalf("Failed to create dbt_project.yml file %v\n", err)
	}
	defer o.Close()
	cd.ProjectName = pn
	err = tmpl.Execute(o, cd)
	if err != nil {
		log.Fatalf("Failed to execute dbt_project.yml template %v\n", err)
	}
	gi := path.Join(bd, ".gitignore")
	err = os.WriteFile(gi, gitignore, 0644)
	if err != nil {
		log.Fatalf("Failed to write .gitignore file %v\n", err)
	}
	s := path.Join(bd, "models/staging", cd.Schema)
	err = os.MkdirAll(s, 0755)
	if err != nil {
		return "", err
	}
	return s, nil
}
