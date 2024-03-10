package main

import (
	"log"
	"os"
	"path/filepath"
)

func CleanBuildDir(buildDir string) {
	_, err := os.Stat(buildDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(buildDir, 0755)
		if errDir != nil {
			panic(err)
		}
	}

	err = filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error cleaning build directory: %v", err)
	}
}
