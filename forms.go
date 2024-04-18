package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

type FormResponse struct {
	Confirm              bool
	Warehouse            string
	Username             string
	Account              string
	Database             string
	Schema               string
	Project              string
	Dataset              string
	Path                 string
	BuildDir             string
	GenerateDescriptions bool
	GroqKeyEnvVar        string
	UseDbtProfile        bool
	DbtProfile           string
	DbtProfileOutput     string
	CreateProfile        bool
	ScaffoldProject      bool
	ProjectName          string
	Prefix               string
}

func Forms() (formResponse FormResponse) {
	formResponse = FormResponse{}
	introForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🏁 Welcome to tbd! 🏎️✨").
				Description(fmt.Sprintf(`A sweet and speedy code generator for dbt.
¸.•✴︎•.¸.•✴︎•.¸.•✴︎•. _%s_ .•✴︎•.¸.•✴︎•.¸.•✴︎•.¸
Currently supports *Snowflake*, *BigQuery*, and *DuckDB*.

Generates:
✴︎ YAML sources config
✴︎ SQL staging models
For each table in the designated schema/dataset.

To prepare, make sure you have the following:
✴︎ The name of an existing dbt profile to reference
(Can be found in the profiles.yml file)
*_OR_*
✴︎ The necessary connection details for your warehouse

_Authentication must be handled via SSO._
_For security, we don't support password auth._

Platform-specific requirements:
✴︎ _Snowflake_: externalbrowser auth
✴︎ _BigQuery_: gcloud CLI installed and authed
✴︎ _DuckDB_: none if using a local db
`, Version)),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("🤖 Experimental: LLM Generation 🦙✨").
				Description(`*_Optional_* LLM-powered alpha features.

Currently generates: 
✴︎ column _descriptions_
✴︎ relevant _tests_
via the Groq API.

You'll need:
✴︎ A Groq API key
✴︎ Key stored in env var`),
			huh.NewConfirm().Affirmative("Sure!").Negative("Nope").
				Title("Do you want to generate column descriptions and tests via LLM?").
				Value(&formResponse.GenerateDescriptions),
		),
		huh.NewGroup(
			huh.NewConfirm().Affirmative("Yes!").Negative("Nah").
				Title("Do you have a dbt profile you'd like to connect with?\n(you can enter your credentials manually if not)").
				Value(&formResponse.UseDbtProfile),
			huh.NewConfirm().Affirmative("Yeah!").Negative("Nope").
				Title("Would you like to scaffold a basic dbt project into the output directory?").
				Value(&formResponse.ScaffoldProject),
			huh.NewInput().
				Title("What prefix do you want to use for your staging files?").
				Value(&formResponse.Prefix).
				Placeholder("stg"),
		),
	)
	projectNameForm := huh.NewForm(
		huh.NewGroup(huh.NewInput().
			Title("What is the name of your dbt project?").
			Value(&formResponse.ProjectName).
			Placeholder("gondor_patrol_analytics"),
		))
	profileCreateForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Affirmative("Yes, pls").Negative("No, thx").
				Title("Would you like to generate a profiles.yml file from the info you provide next?").
				Value(&formResponse.CreateProfile),
		))
	dbtForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is the dbt profile name you'd like to use?").
				Value(&formResponse.DbtProfile).
				Placeholder("snowflake_sandbox"),
			huh.NewInput().
				Title("Which 'output' in that profile do you want to use?").
				Value(&formResponse.DbtProfileOutput).
				Placeholder("dev"),
			huh.NewInput().
				Title("What schema/dataset do you want to generate?").
				Value(&formResponse.Schema),
		),
	)
	warehouseForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your warehouse.").
				Options(
					huh.NewOption("Snowflake", "snowflake"),
					huh.NewOption("BigQuery", "bigquery"),
					huh.NewOption("DuckDB", "duckdb"),
				).
				Value(&formResponse.Warehouse),
		),
	)
	snowflakeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is your username?").
				Value(&formResponse.Username).Placeholder("aragorn@dunedain.king"),

			huh.NewInput().
				Title("What is your Snowflake account id?").
				Value(&formResponse.Account).Placeholder("elfstone-consulting.us-west-1"),

			huh.NewInput().
				Title("What is the schema you want to generate?").
				Value(&formResponse.Schema).Placeholder("minas-tirith"),

			huh.NewInput().
				Title("What database is that schema in?").
				Value(&formResponse.Database).Placeholder("gondor"),
		),
	)
	bigqueryForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("What is your GCP project's id?").
				Value(&formResponse.Project).Placeholder("legolas_inc"),
			huh.NewInput().Title("What is the dataset you want to generate?").
				Value(&formResponse.Dataset).Placeholder("mirkwood"),
		),
	)
	duckdbForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title(`What is the path to your DuckDB database?
Relative to pwd e.g. if db is in this dir -> cool_ducks.db`).
				Value(&formResponse.Path).Placeholder("/path/to/duckdb.db"),
			huh.NewInput().Title("What is the DuckDB database you want to generate?").
				Value(&formResponse.Database).Placeholder("duckdb"),
			huh.NewInput().Title("What is the schema you want to generate?").
				Value(&formResponse.Schema).Placeholder("raw"),
		),
	)
	llmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What env var holds your Groq key?").
				Placeholder("GROQ_API_KEY").
				Value(&formResponse.GroqKeyEnvVar),
		),
	)
	dirForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🚧🚨 Choose your build directory carefully! 🚨🚧").
				Description(`Choose a _new_ or _empty_ directory.
If you choose an existing, populated directory 
tbd will _intentionally error out_.`),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What directory do you want to build into?").
				Value(&formResponse.BuildDir).
				Placeholder("build"),
		),
	)
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Affirmative("Let's go!").Negative("Nevermind").
				Title("🚦Are you ready to do this thing?🚦").
				Value(&formResponse.Confirm),
		),
	)
	introForm.WithTheme(huh.ThemeCatppuccin())
	profileCreateForm.WithTheme(huh.ThemeCatppuccin())
	projectNameForm.WithTheme(huh.ThemeCatppuccin())
	dbtForm.WithTheme(huh.ThemeCatppuccin())
	warehouseForm.WithTheme(huh.ThemeCatppuccin())
	snowflakeForm.WithTheme(huh.ThemeCatppuccin())
	bigqueryForm.WithTheme(huh.ThemeCatppuccin())
	duckdbForm.WithTheme(huh.ThemeCatppuccin())
	llmForm.WithTheme(huh.ThemeCatppuccin())
	dirForm.WithTheme(huh.ThemeCatppuccin())
	confirmForm.WithTheme(huh.ThemeCatppuccin())
	err := introForm.Run()
	if err != nil {
		log.Fatalf("Error running intro form %v\n", err)
	}
	if formResponse.UseDbtProfile {
		err = dbtForm.Run()
		if err != nil {
			log.Fatalf("Error running dbt form %v\n", err)
		}
	} else {
		err = profileCreateForm.Run()
		if err != nil {
			log.Fatalf("Error running profile create form %v\n", err)
		}
		if formResponse.ScaffoldProject {
			err = projectNameForm.Run()
			if err != nil {
				log.Fatalf("Error running project name form %v\n", err)
			}
		}
		err = warehouseForm.Run()
		if err != nil {
			log.Fatalf("Error running warehouse form %v\n", err)
		}
		switch formResponse.Warehouse {
		case "snowflake":
			err = snowflakeForm.Run()
			if err != nil {
				log.Fatalf("Error running snowflake form %v\n", err)
			}
		case "bigquery":
			{
				err = bigqueryForm.Run()
				if err != nil {
					log.Fatalf("Error running bigquery form %v\n", err)
				}
			}
		case "duckdb":
			{
				err = duckdbForm.Run()
				if err != nil {
					log.Fatalf("Error running duckdb form %v\n", err)
				}
			}
		}
	}
	if formResponse.GenerateDescriptions {
		err = llmForm.Run()
		if err != nil {
			log.Fatalf("Error running LLM features form %v\n", err)
		}
	}
	err = dirForm.Run()
	if err != nil {
		log.Fatalf("Error running build directory form %v\n", err)
	}
	err = confirmForm.Run()
	if err != nil {
		log.Fatalf("Error running confirmation form %v\n", err)
	}
	return formResponse
}
