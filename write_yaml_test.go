package main

import (
	"os"
	"testing"

	"github.com/gwenwindflower/tbd/shared"
)

func TestWriteYAML(t *testing.T) {
	tablesFixture := shared.SourceTables{
		SourceTables: []shared.SourceTable{
			{
				Name: "table1",
				Columns: []shared.Column{
					{
						Name:        "column1",
						Description: "column1 description",
						DataType:    "int",
						Tests:       []string{},
					},
				},
				DataTypeGroups: map[string][]shared.Column{
					"int": {
						shared.Column{
							Name:        "column1",
							Description: "column1 description",
							DataType:    "int",
							Tests:       []string{},
						},
					},
				},
			},
		},
	}
	buildDir := "testWriteYAML"
	PrepBuildDir(buildDir)
	WriteYAML(tablesFixture, buildDir)
	_, err := os.Stat(buildDir + "/_sources.yml")
	if os.IsNotExist(err) {
		t.Errorf("WriteYAML did not create the file")
	}
	err = os.RemoveAll(buildDir)
	if err != nil {
		t.Errorf("Failed to clean up test directory")
	}
}
