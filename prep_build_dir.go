package main

import (
	"log"
	"os"
)

func PrepBuildDir(buildDir string) {
	_, err := os.Stat(buildDir)
	if os.IsNotExist(err) {
		dirErr := os.MkdirAll(buildDir, 0755)
		if dirErr != nil {
			log.Fatalf("Failed to create directory %v", dirErr)
		}
	}
}
