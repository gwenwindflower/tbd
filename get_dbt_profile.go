package main

import (
	"fmt"
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
		Threads       int    `yaml:"threads"`
	} `yaml:"outputs"`
}

func GetDbtProfile(dbtProfile string) (*DbtProfile, error) {
	paths := []string{
		filepath.Join(".", "profiles.yml"),
		filepath.Join(os.Getenv("HOME"), ".dbt", "profiles.yml"),
	}
	profileMap := make(map[string]DbtProfile)
	var selectedProfile *DbtProfile
	for _, path := range paths {
		yamlFile, err := os.ReadFile(path)
		if err == nil {
			if err := yaml.Unmarshal(yamlFile, &profileMap); err != nil {
				log.Fatalf("Could not unmarshal dbt profile: %v\n", err)
			}

			if profile, ok := profileMap[dbtProfile]; ok {
				selectedProfile = &profile
			}
		}
	}
	if selectedProfile != nil {
		return selectedProfile, nil
	} else {
		return nil, fmt.Errorf("no profile named %s", dbtProfile)
	}
}
