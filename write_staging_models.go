package main

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"
)

//go:embed *.sql
var stagingTemplate embed.FS

func WriteStagingModels(tables SourceTables, buildDir string) {
	var wg sync.WaitGroup

	for _, table := range tables.SourceTables {
		wg.Add(1)
		go func(table SourceTable) {
			defer wg.Done()

			tmpl := template.New("staging_template.sql").Funcs(template.FuncMap{"lower": strings.ToLower})
			tmpl, err := tmpl.ParseFS(stagingTemplate, "staging_template.sql")
			if err != nil {
				panic(err)
			}

			_, err = os.Stat(buildDir)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(buildDir, 0755)
				if errDir != nil {
					panic(err)
				}
			}

			filename := fmt.Sprintf(buildDir + "/stg_" + strings.ToLower(table.Name) + ".sql")
			outputFile, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()

			err = tmpl.Execute(outputFile, table)
			if err != nil {
				panic(err)
			}
		}(table)
	}
	wg.Wait()
}
