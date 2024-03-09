package main

import (
	"fmt"
	"log"
	"strings"

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
			huh.NewNote().
				Title("Welcome to tbd!🏎️✨").
				Description(`A fast and friendly code generator for dbt.
We will generate sources YAML config and SQL staging models for all the tables in the schema you specify.
To prepare, make sure you have the following:
✴︎ *_Username_* (e.g. aragorn@dunedain.king)
✴︎ *_Account ID_* (e.g. elfstone-consulting.us-west-1)
✴︎ *_Schema_* you want to generate (e.g. minas-tirith)
✴︎ *_Database_* that schema is in (e.g. gondor)
Authentication will be handled via SSO in the web browser.
For security, we don't currently support password-based authentication.`),
		),
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
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you ready to go?").
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
		if warehouse == "snowflake" {
			dbAccount = strings.ToUpper(dbAccount)
			dbUsername = strings.ToUpper(dbUsername)
			dbSchema = strings.ToUpper(dbSchema)
			dbDatabase = strings.ToUpper(dbDatabase)
		}
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
		CleanBuildDir("build")
		WriteYAML(tables)
		WriteStagingModels(tables)
	}
}
