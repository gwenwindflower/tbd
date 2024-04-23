package internal

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/gwenwindflower/tbd/shared"
)

//go:embed *.sql
var stagingTemplate embed.FS

func WriteStagingModels(ts shared.SourceTables, buildDir string, prefix string) {
	var wg sync.WaitGroup

	for _, t := range ts.SourceTables {
		wg.Add(1)
		go func(t shared.SourceTable) {
			defer wg.Done()

			tmpl := template.New("staging_template.sql").Funcs(template.FuncMap{"lower": strings.ToLower})
			tmpl, err := tmpl.ParseFS(stagingTemplate, "staging_template.sql")
			if err != nil {
				log.Fatalf("Failed to parse template %v\n", err)
			}

			filename := fmt.Sprintf(buildDir + "/" + prefix + "_" + strings.ToLower(t.Name) + ".sql")
			outputFile, err := os.Create(filename)
			if err != nil {
				log.Fatalf("Failed to create file %v\n", err)
			}
			defer outputFile.Close()

			err = tmpl.Execute(outputFile, t)
			if err != nil {
				log.Fatalf("Failed to execute template %v\n", err)
			}
		}(t)
	}
	wg.Wait()
}
