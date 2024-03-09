package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func WriteYAML(tables SourceTables) {
	yamlData, err := yaml.Marshal(tables)
	if err != nil {
		log.Fatal(err)
	}

	writeError := os.WriteFile("build/_sources.yml", yamlData, 0644)
	if writeError != nil {
		log.Fatalf("Failed to write file %v", writeError)
	}
}
