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
tbd will generate source YAML config and SQL staging models for all the tables in the schema you specify.
To prepare, make sure you have the following:
* An existing dbt profile.yml file to reference
OR
✴︎ *_Username_* (e.g. aragorn@dunedain.king)
✴︎ *_Account ID_* (e.g. elfstone-consulting.us-west-1)
✴︎ *_Schema_* you want to generate (e.g. minas-tirith)
✴︎ *_Database_* that schema is in (e.g. gondor)
Authentication will be handled via SSO in the web browser.
For security, we don't currently support password-based authentication.`),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("⚠️ Experimental Feature: LLM Generation 🦙✨").
				Description(`I'm currently exploring *_optional_* LLM-powered alpha features.
At present this is limited to generating column descriptions and inferring tests via Groq.
You'll need:
✴︎ A Groq API key stored in an environment variable.`),
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
		),
	)
	manual_form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your warehouse.").
				Options(
					huh.NewOption("Snowflake", "snowflake"),
				).
				Value(&formResponse.Warehouse),

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
	llm_form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is the name of the env var you've stored your Groq API key in?").
				Placeholder("GROQ_API_KEY").
				Value(&formResponse.GroqKeyEnvVar),
		),
	)
	dir_form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🚧🚨 Choose your build directory carefully! 🚨🚧").
				Description(`I highly recommend choosing a new or empty directory to build into.
If you use an existing directory, tbd will overwrite any existing files with the same name.`),
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
	manual_form.WithTheme(huh.ThemeCatppuccin())
	llm_form.WithTheme(huh.ThemeCatppuccin())
	dir_form.WithTheme(huh.ThemeCatppuccin())
	confirm_form.WithTheme(huh.ThemeCatppuccin())
	err := intro_form.Run()
	if err != nil {
		log.Fatal(err)
	}
	if formResponse.UseDbtProfile {
		err = dbt_form.Run()
	} else {
		err = manual_form.Run()
	}
	if err != nil {
		log.Fatal(err)
	}
	if formResponse.GenerateDescriptions {
		err = llm_form.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = dir_form.Run()
	if err != nil {
		log.Fatal(err)
	}
	err = confirm_form.Run()
	if err != nil {
		log.Fatal(err)
	}
	return formResponse
}
