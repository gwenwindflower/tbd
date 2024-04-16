package main

import (
	"errors"
	"log"
	"os"
)

func PrepBuildDir(bd string) error {
	_, err := os.Stat(bd)
	if os.IsNotExist(err) {
		dirErr := os.MkdirAll(bd, 0755)
		if dirErr != nil {
			return dirErr
		}
	} else if err == nil {
		files, err := os.ReadDir(bd)
		if err != nil {
			log.Fatalf("Failed to check build target directory %v", err)
		}
		if len(files) == 0 {
			return nil
		} else {
			return errors.New("build directory is not empty")
		}
	} else {
		return err
	}
	return nil
}
