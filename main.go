package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	_ "github.com/snowflakedb/gosnowflake"
)

// Connection details for Snowflake
var (
	confirm    bool
	warehouse  string
	dbUsername string
	dbAccount  string
	dbDatabase string
	dbSchema   string
)

// Type definitions for the YAML file
type Column struct {
	Name     string `yaml:"name"`
	DataType string `yaml:"data_type"`
}

type SourceTable struct {
	Name    string   `yaml:"name"`
	Columns []Column `yaml:"columns"`
}

type SourceTables struct {
	SourceTables []SourceTable `yaml:"sources"`
}

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your warehouse.").
				Options(
					huh.NewOption("Snowflake", "snowflake"),
				).
				Value(&warehouse),

			huh.NewInput().
				Title("What is your username?").
				Value(&dbUsername),

			huh.NewInput().
				Title("What is your Snowflake account id?").
				Value(&dbAccount),

			huh.NewInput().
				Title("What is the schema you want to generate?").
				Value(&dbSchema),

			huh.NewInput().
				Title("What database is that schema in?").
				Value(&dbDatabase),

			huh.NewConfirm().
				Title("Are you ready to proceed?").
				Value(&confirm),
		),
	)
	form.WithTheme(huh.ThemeCatppuccin())
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	if confirm {
		databaseType := warehouse
		connStr := fmt.Sprintf("%s@%s/%s/%s?authenticator=externalbrowser", dbUsername, dbAccount, dbDatabase, dbSchema)
		ctx, db, err := ConnectToDB(connStr, databaseType)
		if err != nil {
			log.Fatal(err)
		}
		tables, err := GetTables(db, ctx)
		if err != nil {
			log.Fatal(err)
		}
		PutColumnsOnTables(db, ctx, tables)
		WriteYAML(tables)

		WriteStagingModels(tables)
	}
}
