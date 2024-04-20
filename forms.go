package main

import (
	"fmt"

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
	DbtProfileName       string
	DbtProfileOutput     string
	CreateProfile        bool
	ScaffoldProject      bool
	ProjectName          string
	Prefix               string
}

var not_empty = func(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("cannot be empty, please enter a value")
	}
	return nil
}

func getProfileOptions(ps DbtProfiles) []huh.Option[string] {
	var po []huh.Option[string]
	for k := range ps {
		po = append(po, huh.Option[string]{
			Key:   k,
			Value: k,
		})
	}
	return po
}

func Forms(ps DbtProfiles) (FormResponse, error) {
	dfr := FormResponse{
		BuildDir:      "build",
		GroqKeyEnvVar: "GROQ_API_KEY",
		Prefix:        "stg",
	}
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🏁 Welcome to tbd! 🏎️✨").
				Description(fmt.Sprintf(`A sweet and speedy code generator for dbt.
¸.•✴︎•.¸.•✴︎•.¸.•✴︎•. _%s_ .•✴︎•.¸.•✴︎•.¸.•✴︎•.¸
To prepare, make sure you have the following:

✴︎ The name of an *_existing dbt profile_* to reference
*_OR_*
✴︎ The necessary *_connection details_* for your warehouse

_See README for warehouse-specific requirements_
`, Version)),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you have a dbt profile you'd like to connect with?\n(you can enter your credentials manually if not)").
				Value(&dfr.UseDbtProfile),
			huh.NewConfirm().
				Title("Would you like to scaffold a basic dbt project into the output directory?").
				Value(&dfr.ScaffoldProject),
			huh.NewInput().
				Title("What prefix do you want to use for your staging files?").
				Value(&dfr.Prefix).
				Placeholder("stg").
				Validate(not_empty),
		),

		huh.NewGroup(huh.NewInput().
			Title("What is the name of your dbt project?").
			Value(&dfr.ProjectName).
			Placeholder("gondor_patrol_analytics").
			Validate(not_empty),
		).WithHideFunc(func() bool {
			return !dfr.ScaffoldProject
		}),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to generate a profiles.yml file dfrom the info you provide next?").
				Value(&dfr.CreateProfile),
		).WithHideFunc(func() bool {
			return dfr.UseDbtProfile
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a dbt profile:").
				Options(getProfileOptions(ps)...).
				Value(&dfr.DbtProfileName),
			huh.NewInput().
				Title("Which 'output' in that profile do you want to use?").
				Value(&dfr.DbtProfileOutput).
				Placeholder("dev").
				Validate(not_empty),
			huh.NewInput().
				Title("What schema/dataset do you want to generate?").
				Value(&dfr.Schema).
				Placeholder("raw").
				Validate(not_empty),
			huh.NewInput().
				Title("What project/database is that schema/dataset in?").
				Value(&dfr.Schema).
				Placeholder("jaffle_shop").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return !dfr.UseDbtProfile
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your warehouse.").
				Options(
					huh.NewOption("Snowflake", "snowflake"),
					huh.NewOption("BigQuery", "bigquery"),
					huh.NewOption("DuckDB", "duckdb"),
				).
				Value(&dfr.Warehouse),
		).WithHideFunc(func() bool {
			return dfr.UseDbtProfile
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What is your username?").
				Value(&dfr.Username).
				Placeholder("aragorn@dunedain.king").
				Validate(not_empty),
			huh.NewInput().
				Title("What is your Snowflake account id?").
				Value(&dfr.Account).
				Placeholder("elfstone-consulting.us-west-1").
				Validate(not_empty),
			huh.NewInput().
				Title("What is the schema you want to generate?").
				Value(&dfr.Schema).
				Placeholder("minas-tirith").
				Validate(not_empty),
			huh.NewInput().
				Title("What database is that schema in?").
				Value(&dfr.Database).
				Placeholder("gondor").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "snowflake"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What is your GCP project's id?").
				Value(&dfr.Project).
				Placeholder("legolas_inc").
				Validate(not_empty),
			huh.NewInput().
				Title("What is the dataset you want to generate?").
				Value(&dfr.Dataset).
				Placeholder("mirkwood").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "bigquery"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title(`What is the path to your DuckDB database?
Relative to pwd e.g. if db is in this dir -> cool_ducks.db`).
				Value(&dfr.Path).
				Placeholder("/path/to/duckdb.db").
				Validate(not_empty),
			huh.NewInput().
				Title("What is the DuckDB database you want to generate?").
				Value(&dfr.Database).
				Placeholder("duckdb").
				Validate(not_empty),
			huh.NewInput().
				Title("What is the schema you want to generate?").
				Value(&dfr.Schema).
				Placeholder("raw").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "duckdb"
		}),

		huh.NewGroup(
			huh.NewNote().
				Title("🤖 Experimental: LLM Generation 🦙✨").
				Description(`*_Optional_* LLM-powered alpha features, powered by Groq.

Currently generates: 
✴︎ column _descriptions_
✴︎ relevant _tests_

You'll need:
✴︎ A Groq API key in an env var`),
			huh.NewConfirm().
				Title("Do you want to generate column descriptions and tests via LLM?").
				Value(&dfr.GenerateDescriptions),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("What env var holds your Groq key?").
				Placeholder("GROQ_API_KEY").
				Value(&dfr.GroqKeyEnvVar).
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return !dfr.GenerateDescriptions
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What directory do you want to build into?\n Must be new or empty.").
				Value(&dfr.BuildDir).
				Placeholder("build").
				Validate(not_empty),
			huh.NewConfirm().
				Title("🚦Are you ready to do this thing?🚦").
				Value(&dfr.Confirm),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	if err != nil {
		return dfr, err
	}
	return dfr, nil
}
