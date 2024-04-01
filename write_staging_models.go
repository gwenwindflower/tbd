package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"tbd/shared"
	"text/template"
)

//go:embed *.sql
var stagingTemplate embed.FS

func WriteStagingModels(tables shared.SourceTables, buildDir string) {
	var wg sync.WaitGroup

	for _, table := range tables.SourceTables {
		wg.Add(1)
		go func(table shared.SourceTable) {
			defer wg.Done()

			tmpl := template.New("staging_template.sql").Funcs(template.FuncMap{"lower": strings.ToLower})
			tmpl, err := tmpl.ParseFS(stagingTemplate, "staging_template.sql")
			if err != nil {
				log.Fatalf("Failed to parse template %v\n", err)
			}

			filename := fmt.Sprintf(buildDir + "/stg_" + strings.ToLower(table.Name) + ".sql")
			outputFile, err := os.Create(filename)
			if err != nil {
				log.Fatalf("Failed to create file %v\n", err)
			}
			defer outputFile.Close()

			err = tmpl.Execute(outputFile, table)
			if err != nil {
				log.Fatalf("Failed to execute template %v\n", err)
			}
		}(table)
	}
	wg.Wait()
}
