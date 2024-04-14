package main

import (
	"os"
	"testing"
)

func TestGetDbtProfile(t *testing.T) {
	CreateTempDbtProfile(t)
	defer os.RemoveAll(os.Getenv("HOME"))
	defer os.Unsetenv("HOME")
	// Profile exists
	profile, err := GetDbtProfile("elf")
	if err != nil {
		t.Errorf("GetDbtProfile returned an error for an existing profile: %v", err)
	}
	if profile.Target != "dev" {
		t.Errorf("Expected target 'dev', got '%s'", profile.Target)
	}

	// Profile does not exist
	profile, err = GetDbtProfile("dunedain")
	if err == nil {
		t.Error("GetDbtProfile did not return an error for a non-existing profile")
	}
	if profile != nil {
		t.Error("GetDbtProfile returned a non-nil profile for a non-existing profile")
	}
}
