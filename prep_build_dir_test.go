package main

import (
	"os"
	"testing"
)

func TestPrepBuildDir(t *testing.T) {
	buildDir := "testPrepBuildDir"
	PrepBuildDir(buildDir)
	_, err := os.Stat(buildDir)
	if os.IsNotExist(err) {
		t.Errorf("PrepBuildDir did not create the directory")
	}
	os.Remove(buildDir)
}
