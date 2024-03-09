package main

import (
	"log"
	"os"
	"path/filepath"
)

func CleanBuildDir(buildDirPath string) {
	err := filepath.Walk(buildDirPath, func(path string, info os.FileInfo, err error) error {
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
