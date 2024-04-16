package main

import (
	"log"
	"os"
	"path"
	"text/template"

	"github.com/gwenwindflower/tbd/shared"
)

func WriteProfile(cd shared.ConnectionDetails, bd string) {
	pt := `
{{.ConnType}}:
  target: dev
  outputs:
    dev:
      type: {{.ConnType}}
      {{- if eq .ConnType "snowflake"}}
      auth: externalbrowser
      {{- end}}
      {{- if eq .ConnType "bigquery"}}
      method: oauth
      {{- end}}
      {{- if .Account}}
      account: {{.Account}}
      {{- end}}
      {{- if .Username}}
      user: {{.Username}}
      {{- end}}
      {{- if .Database}}
      database: {{.Database}}
      {{- end}}
      {{- if .Project}}
      project: {{.Project}}
      {{- end}}
      {{- if .Schema}}
      schema: {{.Schema}}
      {{- end}}
      {{- if .Dataset}}
      dataset: {{.Dataset}}
      {{- end}}
      {{- if .Path}}
      path: {{.Path}}
      {{- end}}
      threads: 8
`
	tmpl, err := template.New("profiles.yml").Parse(pt)
	if err != nil {
		log.Fatalf("Failed to parse template %v\n", err)
	}
	p := path.Join(bd, "profiles.yml")
	o, err := os.Create(p)
	if err != nil {
		log.Fatalf("Failed to create profiles.yml file %v\n", err)
	}
	defer o.Close()
	err = tmpl.Execute(o, cd)
	if err != nil {
		log.Fatalf("Failed to execute template %v\n", err)
	}
}
