package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func WriteStagingModels(tables SourceTables) {
	for _, table := range tables.SourceTables {
		go func(table SourceTable) {
			tmpl, err := template.ParseFiles("./template.sql")
			if err != nil {
				log.Fatal(err)
			}

			filename := fmt.Sprintf("build/stg_" + strings.ToLower(table.Name) + ".sql")
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
}
