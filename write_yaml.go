package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func WriteYAML(tables SourceTables, buildDir string) {
	yamlData, err := yaml.Marshal(tables)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Refactor this to a function we can reuse, it's called in multiple places
	_, err = os.Stat(buildDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(buildDir, 0755)
		if errDir != nil {
			panic(err)
		}
	}

	writeError := os.WriteFile(buildDir+"/_sources.yml", yamlData, 0644)
	if writeError != nil {
		log.Fatalf("Failed to write file %v", writeError)
	}
}
