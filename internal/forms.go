package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
)

type FormResponse struct {
	Password             string
	Host                 string
	BuildDir             string
	SslMode              string
	Database             string
	Schema               string
	Project              string
	Dataset              string
	ProjectName          string
	Warehouse            string
	Account              string
	LlmKeyEnvVar         string
	DbtProfileOutput     string
	DbtProfileName       string
	Path                 string
	Port                 string
	TokenEnvVar          string
	HttpPath             string
	Username             string
	Prefix               string
	Llm                  string
	GenerateDescriptions bool
	ScaffoldProject      bool
	CreateProfile        bool
	UseDbtProfile        bool
	Confirm              bool
}

func notEmpty(s string) error {
	s = strings.TrimSpace(s)
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
		BuildDir:     "build",
		LlmKeyEnvVar: "OPENAI_API_KEY",
		Prefix:       "stg",
		Port:         "5432",
		TokenEnvVar:  "DATABRICKS_TOKEN",
	}
	pinkUnderline := color.New(color.FgMagenta).Add(color.Bold, color.Underline).SprintFunc()
	greenBold := color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	yellowItalic := color.New(color.FgHiYellow).Add(color.Italic).SprintFunc()
	greenBoldItalic := color.New(color.FgHiGreen).Add(color.Bold).SprintFunc()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🏁 Welcome to tbd! 🏎️✨").
				Description(fmt.Sprintf(`A sweet and speedy code generator for dbt.
¸.•✴︎•.¸.•✴︎•.¸.•✴︎•. _%s_ .•✴︎•.¸.•✴︎•.¸.•✴︎•.¸
To prepare, make sure you have the following:

✴︎ The name of an %s to reference
*_OR_*
✴︎ The necessary %s for your warehouse

_See_ %s _for warehouse-specific requirements_:
https://github.com/gwenwindflower/tbd
`, greenBold(VERSION), pinkUnderline("existing dbt profile"), pinkUnderline("connection details"), greenBold("README"))),
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
				Validate(notEmpty),
		),

		huh.NewGroup(huh.NewInput().
			Title("What is the *name* of your dbt project?").
			Value(&dfr.ProjectName).
			Placeholder("rivendell").
			Validate(notEmpty),
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
				Validate(notEmpty),
			huh.NewInput().
				Title("What *schema* do you want to generate?").
				Value(&dfr.Schema).
				Placeholder("raw").
				Validate(notEmpty),
			huh.NewInput().
				Title("What *database/project/catalog* is that schema in?").
				Value(&dfr.Database).
				Placeholder("jaffle_shop").
				Validate(notEmpty),
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
					huh.NewOption("Postgres", "postgres"),
					huh.NewOption("Databricks", "databricks"),
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
				Validate(notEmpty),
			huh.NewInput().
				Title("What is your Snowflake account id?").
				Value(&dfr.Account).
				Placeholder("elfstone-consulting.us-west-1").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is the *schema* you want to generate?").
				Value(&dfr.Schema).
				Placeholder("minas-tirith").
				Validate(notEmpty),
			huh.NewInput().
				Title("What *database* is that schema in?").
				Value(&dfr.Database).
				Placeholder("gondor").
				Validate(notEmpty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "snowflake"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What GCP *project id* do you want to generate?").
				Value(&dfr.Project).
				Placeholder("legolas_inc").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is the *dataset* you want to generate?").
				Value(&dfr.Dataset).
				Placeholder("mirkwood").
				Validate(notEmpty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "bigquery"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title(`What is the *path* to your DuckDB database?
Relative to pwd e.g. if db is in this dir -> cool_ducks.db`).
				Value(&dfr.Path).
				Placeholder("/path/to/duckdb.db").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is the *database* you want to generate?").
				Value(&dfr.Database).
				Placeholder("gimli_corp").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is the *schema* you want to generate?").
				Value(&dfr.Schema).
				Placeholder("moria").
				Validate(notEmpty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "duckdb"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What is your Postgres *host*?").
				Value(&dfr.Host).
				Validate(notEmpty),
			huh.NewInput().
				Title("What is your Postgres *port*?").
				Value(&dfr.Port).
				Validate(func(s string) error {
					port, err := strconv.Atoi(s)
					if err != nil || port < 1000 || port > 9999 {
						return errors.New("port must be a 4-digit number")
					}
					return nil
				}),
			huh.NewInput().
				Title("What is your Postgres *username*?").
				Value(&dfr.Username).
				Placeholder("galadriel").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is your Postgres *password*?").
				Value(&dfr.Password).
				Validate(notEmpty).
				EchoMode(huh.EchoModePassword),
			huh.NewInput().
				Title("What is the *database* you want to generate?").
				Value(&dfr.Database).
				Placeholder("lothlorien").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is the *schema* you want to generate?").
				Value(&dfr.Schema).
				Placeholder("mallorn_trees").
				Validate(notEmpty),
			huh.NewSelect[string]().
				Title("What ssl mode do you want to use?").
				Value(&dfr.SslMode).
				Options(
					huh.NewOption("Disable", "disable"),
					huh.NewOption("Require", "require"),
					huh.NewOption("Verify-ca", "verify-ca"),
					huh.NewOption("Verify-full", "verify-full"),
					huh.NewOption("Prefer", "prefer"),
					huh.NewOption("Allow", "allow")).
				Validate(notEmpty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "postgres"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What is your Databricks *host*?").
				Value(&dfr.Host).
				Placeholder("dbc-12345.cloud.databricks.com").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is your warehouse's *HTTP path*?").
				Value(&dfr.HttpPath).
				Placeholder("/sql/1.0/warehouses/12345").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is your Databricks *username*?").
				Value(&dfr.Username).
				Placeholder("arwen").
				Validate(notEmpty),
			huh.NewInput().
				Title("What env var holds your Databricks *Personal Access Token*?").
				Value(&dfr.TokenEnvVar).
				Validate(notEmpty),
			huh.NewInput().
				Title("What is the *catalog* you want to generate?").
				Value(&dfr.Database).
				Placeholder("rivendell").
				Validate(notEmpty),
			huh.NewInput().
				Title("What is the *schema* you want to generate?").
				Value(&dfr.Schema).
				Placeholder("evenstar").
				Validate(notEmpty),
		).WithHideFunc(func() bool {
			return dfr.Warehouse != "databricks"
		}),

		huh.NewGroup(
			huh.NewNote().
				Title(fmt.Sprintf("🤖 %s LLM generation 🦙✨", yellowItalic("Optional"))).
				Description(fmt.Sprintf(`Infers:
✴︎ column %s
✴︎ relevant %s

_Requires an_ %s _stored in an env var_.`, pinkUnderline("descriptions"), pinkUnderline("tests"), greenBoldItalic("API key"))),
			huh.NewConfirm().
				Title("Do you want to infer descriptions and tests?").
				Value(&dfr.GenerateDescriptions),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your LLM provider:").
				Options(
					huh.NewOption("OpenAI", "openai"),
					huh.NewOption("Groq", "groq"),
					huh.NewOption("Anthropic", "anthropic"),
				).Value(&dfr.Llm),
			huh.NewInput().
				Title("What env var holds your LLM API key?").
				Placeholder("OPENAI_API_KEY").
				Value(&dfr.LlmKeyEnvVar).
				Validate(notEmpty),
		).WithHideFunc(func() bool {
			return !dfr.GenerateDescriptions
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("What directory do you want to build into?\n Must be new or empty.").
				Value(&dfr.BuildDir).
				Placeholder("build").
				Validate(notEmpty),
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
