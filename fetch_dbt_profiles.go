package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type DbtProfile struct {
	Target  string `yaml:"target"`
	Outputs map[string]struct {
		ConnType      string `yaml:"type"`
		Account       string `yaml:"account"`
		User          string `yaml:"user"`
		Role          string `yaml:"role"`
		Authenticator string `yaml:"authenticator"`
		Database      string `yaml:"database"`
		Schema        string `yaml:"schema"`
		Project       string `yaml:"project"`
		Dataset       string `yaml:"dataset"`
		Path          string `yaml:"path"`
		Threads       int    `yaml:"threads"`
	} `yaml:"outputs"`
}

type DbtProfiles map[string]DbtProfile

func FetchDbtProfiles() (DbtProfiles, error) {
	paths := []string{
		filepath.Join(".", "profiles.yml"),
		filepath.Join(os.Getenv("HOME"), ".dbt", "profiles.yml"),
	}
	ps := DbtProfiles{}
	for _, path := range paths {
		pf := DbtProfiles{}
		yf, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err = yaml.Unmarshal(yf, pf); err != nil {
			log.Fatalf("Could not read dbt profile, \nlikely unsupported fields or formatting issues\n please open an issue: %v\n", err)
		}
		for k, v := range pf {
			ps[k] = v
		}
	}
	return ps, nil
}
