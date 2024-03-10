package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
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
	Name           string              `yaml:"name"`
	Columns        []Column            `yaml:"columns"`
	DataTypeGroups map[string][]Column `yaml:"-"`
}

type SourceTables struct {
	SourceTables []SourceTable `yaml:"sources"`
}

func main() {
	useForm := true
	if useForm {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("🏁 Welcome to tbd! 🏎️✨").
					Description(`A sweet and speedy code generator for dbt.
We will generate source YAML config and SQL staging models for all the tables in the schema you specify.
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
	} else {
		warehouse = "snowflake"
		dbUsername = os.Getenv("SNOWFLAKE_SANDBOX_USER")
		dbAccount = os.Getenv("SNOWFLAKE_SANDBOX_ACCOUNT")
		dbDatabase = "ANALYTICS"
		dbSchema = "JAFFLE_SHOP_RAW"
		confirm = true
	}
	if confirm {
		s := spinner.New()
		databaseType := warehouse
		if warehouse == "snowflake" {
			dbAccount = strings.ToUpper(dbAccount)
			dbUsername = strings.ToUpper(dbUsername)
			dbSchema = strings.ToUpper(dbSchema)
			dbDatabase = strings.ToUpper(dbDatabase)
		}
		connStr := fmt.Sprintf("%s@%s/%s/%s?authenticator=externalbrowser", dbUsername, dbAccount, dbDatabase, dbSchema)
		s.Action(func() {
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

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				WriteYAML(tables)
			}()
			go func() {
				defer wg.Done()
				WriteStagingModels(tables)
			}()
			wg.Wait()
		}).Title("🏎️✨ Generating YAML and SQL files...").Run()
		fmt.Println("🏁 Done! Your YAML and SQL files are in the build directory.")
	}
}
