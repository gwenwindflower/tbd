package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
)

func WriteStagingModels(tables SourceTables) {
	var wg sync.WaitGroup

	fmt.Println("Writing staging models...")
	for _, table := range tables.SourceTables {
		wg.Add(1)
		go func(table SourceTable) {
			defer wg.Done()

			tmpl, err := template.ParseFiles("template.sql")
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
	wg.Wait()
}
