package main

import (
	"os"
)

func PrepBuildDir(buildDir string) {
	_, err := os.Stat(buildDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(buildDir, 0755)
		if errDir != nil {
			panic(err)
		}
	}
}
