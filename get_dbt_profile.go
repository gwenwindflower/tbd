package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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
				log.Fatalf("could not unmarshal dbt profile: %v", err)
			}

			if profile, ok := profileMap[dbtProfile]; ok {
				selectedProfile = &profile
			}
		}
	}
	if selectedProfile != nil {
		return selectedProfile, nil
	} else {
		return nil, fmt.Errorf("could not find profile: %s", dbtProfile)
	}
}
