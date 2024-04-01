package main

import (
	"log"
	"os"
)

func PrepBuildDir(buildDir string) {
	_, err := os.Stat(buildDir)
	if err != nil {
		log.Fatalf("Failed to get directory info %v\n", err)
	}
	if os.IsNotExist(err) {
		dirErr := os.MkdirAll(buildDir, 0755)
		if dirErr != nil {
			log.Fatalf("Failed to create directory %v\n", dirErr)
		}
	}
}
