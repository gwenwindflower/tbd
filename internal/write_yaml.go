package internal

import (
	"log"
	"os"

	"github.com/gwenwindflower/tbd/shared"

	"gopkg.in/yaml.v2"
)

func WriteYAML(tables shared.SourceTables, buildDir string) {
	yamlData, err := yaml.Marshal(tables)
	if err != nil {
		log.Fatalf("Failed to marshal data into YAML %v\n", err)
	}

	writeError := os.WriteFile(buildDir+"/_sources.yml", yamlData, 0644)
	if writeError != nil {
		log.Fatalf("Failed to write file %v\n", writeError)
	}
}
