package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gwenwindflower/tbd/shared"
)

func SetConnectionDetails(formResponse FormResponse) shared.ConnectionDetails {
	var cd shared.ConnectionDetails
	if formResponse.UseDbtProfile {
		profile, err := GetDbtProfile(formResponse.DbtProfile)
		if err != nil {
			log.Fatalf("Could not get dbt profile %v\n", err)
		}
		switch profile.Outputs[formResponse.DbtProfileOutput].ConnType {
		case "snowflake":
			{
				cd = shared.ConnectionDetails{
					ConnType: profile.Outputs[formResponse.DbtProfileOutput].ConnType,
					Username: profile.Outputs[formResponse.DbtProfileOutput].User,
					Account:  profile.Outputs[formResponse.DbtProfileOutput].Account,
					Database: profile.Outputs[formResponse.DbtProfileOutput].Database,
					Schema:   formResponse.Schema,
				}
			}
		case "bigquery":
			{
				cd = shared.ConnectionDetails{
					ConnType: profile.Outputs[formResponse.DbtProfileOutput].ConnType,
					Project:  profile.Outputs[formResponse.DbtProfileOutput].Project,
					Dataset:  formResponse.Schema,
				}
			}
		case "duckdb":
			{
				cd = shared.ConnectionDetails{
					ConnType: profile.Outputs[formResponse.DbtProfileOutput].ConnType,
					Path:     profile.Outputs[formResponse.DbtProfileOutput].Path,
					Database: profile.Outputs[formResponse.DbtProfileOutput].Database,
					Schema:   formResponse.Schema,
				}
			}
		default:
			{
				log.Fatalf("Unsupported connection type %v\n", profile.Outputs[formResponse.DbtProfileOutput].ConnType)
			}
		}
	} else {
		switch formResponse.Warehouse {
		case "snowflake":
			{
				cd = shared.ConnectionDetails{
					ConnType: formResponse.Warehouse,
					Username: formResponse.Username,
					Account:  formResponse.Account,
					Schema:   formResponse.Schema,
					Database: formResponse.Database,
				}
			}
		case "bigquery":
			{
				cd = shared.ConnectionDetails{
					ConnType: formResponse.Warehouse,
					Project:  formResponse.Project,
					Dataset:  formResponse.Dataset,
				}
			}
		case "duckdb":
			wd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Could not get working directory: %v\n", err)
			}
			{
				cd = shared.ConnectionDetails{
					ConnType: formResponse.Warehouse,
					Path:     filepath.Join(wd, formResponse.Path),
					Database: formResponse.Database,
					Schema:   formResponse.Schema,
				}
			}
		default:
			{
				log.Fatalf("Unsupported connection type %v\n", formResponse.Warehouse)
			}
		}
	}
	return cd
}
