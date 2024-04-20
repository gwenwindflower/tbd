package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
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
	pinkUnderline := color.New(color.FgMagenta).Add(color.Bold, color.Underline).SprintFunc()
	greenBold := color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	blueBold := color.New(color.FgHiBlue).Add(color.Bold).SprintFunc()
	yellowItalic := color.New(color.FgHiYellow).Add(color.Italic).SprintFunc()
	greenBoldItalic := color.New(color.FgHiGreen).Add(color.Bold).SprintFunc()
	redBold := color.New(color.FgHiRed).Add(color.Bold).SprintFunc()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(fmt.Sprintf("🏁 %s 🏎️✨", blueBold("Welcome to tbd!"))).
				Description(fmt.Sprintf(`A sweet and speedy code generator for dbt.
¸.•✴︎•.¸.•✴︎•.¸.•✴︎•. _%s_ .•✴︎•.¸.•✴︎•.¸.•✴︎•.¸
To prepare, make sure you have the following:

✴︎ The name of an %s to reference
*_OR_*
✴︎ The necessary %s for your warehouse

_See README for warehouse-specific requirements_
https://github.com/gwenwindflower/tbd
`, greenBold(Version), pinkUnderline("existing dbt profile"), pinkUnderline("connection details"))),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you have a *dbt profile* you'd like to connect with?\n(you can enter your credentials manually if not)").
				Value(&dfr.UseDbtProfile),
			huh.NewConfirm().
				Title("Would you like to *scaffold* a basic dbt project?").
				Value(&dfr.ScaffoldProject),
			huh.NewInput().
				Title("What *prefix* for your staging files?").
				Value(&dfr.Prefix).
				Placeholder("stg").
				Validate(not_empty),
		),

		huh.NewGroup(huh.NewInput().
			Title("What is the *name* of your dbt project?").
			Value(&dfr.ProjectName).
			Placeholder("rivendell").
			Validate(not_empty),
		).WithHideFunc(func() bool {
			return !dfr.ScaffoldProject
		}),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to generate a profiles.yml file?\n(from the info you provide next)").
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
				Title("Which *output* in that profile do you want to use?").
				Value(&dfr.DbtProfileOutput).
				Placeholder("dev").
				Validate(not_empty),
			huh.NewInput().
				Title("What *schema* do you want to generate?").
				Value(&dfr.Schema).
				Placeholder("raw").
				Validate(not_empty),
			huh.NewInput().
				Title("What *database* is that schema in?").
				Value(&dfr.Database).
				Placeholder("jaffle_shop").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return !dfr.UseDbtProfile
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your warehouse:").
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
				Title("What is the *schema* you want to generate?").
				Value(&dfr.Schema).
				Placeholder("minas-tirith").
				Validate(not_empty),
			huh.NewInput().
				Title("What *database* is that schema in?").
				Value(&dfr.Database).
				Placeholder("gondor").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "snowflake"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What GCP *project id* do you want to generate?").
				Value(&dfr.Project).
				Placeholder("legolas_inc").
				Validate(not_empty),
			huh.NewInput().
				Title("What is the *dataset* you want to generate?").
				Value(&dfr.Dataset).
				Placeholder("mirkwood").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "bigquery"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title(`What is the *path* to your DuckDB database?
Relative to pwd e.g. if db is in this dir -> cool_ducks.db`).
				Value(&dfr.Path).
				Placeholder("/path/to/duckdb.db").
				Validate(not_empty),
			huh.NewInput().
				Title("What is the *database* you want to generate?").
				Value(&dfr.Database).
				Placeholder("duckdb").
				Validate(not_empty),
			huh.NewInput().
				Title("What is the *schema* you want to generate?").
				Value(&dfr.Schema).
				Placeholder("raw").
				Validate(not_empty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "duckdb"
		}),

		huh.NewGroup(
			huh.NewNote().
				Title(fmt.Sprintf("🤖 %s LLM generation 🦙✨", redBold("Experimental"))).
				Description(fmt.Sprintf(`%s features via Groq.
Currently generates: 
✴︎ column %s
✴︎ relevant %s

Requires a %s stored in an env var
Get one at https://groq.com.`, yellowItalic("Optional"), pinkUnderline("descriptions"), pinkUnderline("tests"), greenBoldItalic("Groq API key"))),
			huh.NewConfirm().
				Title("Do you want to infer descriptions and tests?").
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
