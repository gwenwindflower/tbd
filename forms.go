package main

import (
	"log"

	"github.com/charmbracelet/huh"
)

func Forms() {
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
At present this is limited to generating column descriptions via Groq.
You'll need:
✴︎ A Groq API key stored in an environment variable.`),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you have a dbt profile you'd like to connect with?").
				Value(&useDbtProfile),
		),
	)
	dbt_form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is the dbt profile name you'd like to use?").
				Value(&dbtProfile).
				Placeholder("snowflake_sandbox"),
		),
	)
	main_form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your warehouse.").
				Options(
					huh.NewOption("Snowflake", "snowflake"),
				).
				Value(&warehouse),

			huh.NewInput().
				Title("What is your username?").
				Value(&dbUsername).Placeholder("aragorn@dunedain.king"),

			huh.NewInput().
				Title("What is your Snowflake account id?").
				Value(&dbAccount).Placeholder("elfstone-consulting.us-west-1"),

			huh.NewInput().
				Title("What is the schema you want to generate?").
				Value(&dbSchema).Placeholder("minas-tirith"),

			huh.NewInput().
				Title("What database is that schema in?").
				Value(&dbDatabase).Placeholder("gondor"),

			huh.NewConfirm().
				Title("Do you want to generate column descriptions via LLM?\n⚠️  Experimental ⚠️").
				Value(&generateDescriptions),
		),
	)
	llm_form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is the name of the env var you've stored your Groq API key in?").
				Placeholder("GROQ_API_KEY").
				Value(&groqKeyEnvVar),
		),
	)
	confirm_form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What directory do you want to build into?\n🚧 Name a new or empty directory 🚧").
				Value(&buildDir).Placeholder("build"),

			huh.NewConfirm().
				Title("🚦Are you ready to do this thing?🚦").
				Value(&confirm),
		),
	)
	intro_form.WithTheme(huh.ThemeCatppuccin())
	dbt_form.WithTheme(huh.ThemeCatppuccin())
	main_form.WithTheme(huh.ThemeCatppuccin())
	llm_form.WithTheme(huh.ThemeCatppuccin())
	confirm_form.WithTheme(huh.ThemeCatppuccin())
	err := intro_form.Run()
	if err != nil {
		log.Fatal(err)
	}
	if useDbtProfile {
		err = dbt_form.Run()
	} else {
		err = main_form.Run()
	}
	if err != nil {
		log.Fatal(err)
	}
	if generateDescriptions {
		err = llm_form.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = confirm_form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
