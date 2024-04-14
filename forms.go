package main

import (
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
	BuildDir             string
	GenerateDescriptions bool
	GroqKeyEnvVar        string
	UseDbtProfile        bool
	DbtProfile           string
	DbtProfileOutput     string
}

func Forms() (formResponse FormResponse) {
	formResponse = FormResponse{}
	intro_form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🏁 Welcome to tbd! 🏎️✨").
				Description(`A sweet and speedy code generator for dbt.
Currently supports _Snowflake_ and _BigQuery_.

Generates:
✴︎ YAML sources config
✴︎ SQL staging models
For each table in the designated schema/dataset.

To prepare, make sure you have the following:
✴︎ An existing dbt profile.yml file to reference
*_OR_*
✴︎ The necessary connection details for your warehouse

_Authentication must be handled via SSO._
_For security, we don't support password auth._
Platform-specific requirements:
*_Snowflake_*: externalbrowser auth
*_BigQuery_*: gcloud CLI installed and authed`),
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
		),
	)
	dbt_form := huh.NewForm(
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
	warehouse_form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your warehouse.").
				Options(
					huh.NewOption("Snowflake", "snowflake"),
					huh.NewOption("BigQuery", "bigquery"),
				).
				Value(&formResponse.Warehouse),
		),
	)
	snowflake_form := huh.NewForm(
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
	bigquery_form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("What is your GCP project's id?").
				Value(&formResponse.Project).Placeholder("legolas_inc"),
			huh.NewInput().Title("What is the dataset you want to generate?").
				Value(&formResponse.Dataset).Placeholder("mirkwood"),
		),
	)
	llm_form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What env var holds your Groq key?").
				Placeholder("GROQ_API_KEY").
				Value(&formResponse.GroqKeyEnvVar),
		),
	)
	dir_form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🚧🚨 Choose your build directory carefully! 🚨🚧").
				Description(`_I highly recommend choosing a new or empty directory to build into._
If you use an existing directory,
tbd will overwrite any existing files of the same name.`),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What directory do you want to build into?").
				Value(&formResponse.BuildDir).
				Placeholder("build"),
		),
	)
	confirm_form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Affirmative("Let's go!").Negative("Nevermind").
				Title("🚦Are you ready to do this thing?🚦").
				Value(&formResponse.Confirm),
		),
	)
	intro_form.WithTheme(huh.ThemeCatppuccin())
	dbt_form.WithTheme(huh.ThemeCatppuccin())
	warehouse_form.WithTheme(huh.ThemeCatppuccin())
	snowflake_form.WithTheme(huh.ThemeCatppuccin())
	bigquery_form.WithTheme(huh.ThemeCatppuccin())
	llm_form.WithTheme(huh.ThemeCatppuccin())
	dir_form.WithTheme(huh.ThemeCatppuccin())
	confirm_form.WithTheme(huh.ThemeCatppuccin())
	err := intro_form.Run()
	if err != nil {
		log.Fatalf("Error running intro form %v\n", err)
	}
	if formResponse.UseDbtProfile {
		err = dbt_form.Run()
	} else {
		err = warehouse_form.Run()
		switch formResponse.Warehouse {
		case "snowflake":
			err = snowflake_form.Run()
			if err != nil {
				log.Fatalf("Error running snowflake form %v\n", err)
			}
		case "bigquery":
			err = bigquery_form.Run()
			if err != nil {
				log.Fatalf("Error running bigquery form %v\n", err)
			}
		}
	}
	if err != nil {
		log.Fatalf("Error running connection details form %v\n", err)
	}
	if formResponse.GenerateDescriptions {
		err = llm_form.Run()
		if err != nil {
			log.Fatalf("Error running LLM features form %v\n", err)
		}
	}
	err = dir_form.Run()
	if err != nil {
		log.Fatalf("Error running build directory form %v\n", err)
	}
	err = confirm_form.Run()
	if err != nil {
		log.Fatalf("Error running confirmation form %v\n", err)
	}
	return formResponse
}
