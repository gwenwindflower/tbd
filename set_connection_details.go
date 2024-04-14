package main

import (
	"log"

	"github.com/gwenwindflower/tbd/shared"
)

func SetConnectionDetails(formResponse FormResponse) shared.ConnectionDetails {
	var connectionDetails shared.ConnectionDetails
	if formResponse.UseDbtProfile {
		profile, err := GetDbtProfile(formResponse.DbtProfile)
		if err != nil {
			log.Fatalf("Could not get dbt profile %v\n", err)
		}
		switch profile.Outputs[formResponse.DbtProfileOutput].ConnType {
		case "snowflake":
			{
				connectionDetails = shared.ConnectionDetails{
					ConnType: profile.Outputs[formResponse.DbtProfileOutput].ConnType,
					Username: profile.Outputs[formResponse.DbtProfileOutput].User,
					Account:  profile.Outputs[formResponse.DbtProfileOutput].Account,
					Database: profile.Outputs[formResponse.DbtProfileOutput].Database,
					Schema:   formResponse.Schema,
				}
			}
		case "bigquery":
			{
				connectionDetails = shared.ConnectionDetails{
					ConnType: profile.Outputs[formResponse.DbtProfileOutput].ConnType,
					Project:  profile.Outputs[formResponse.DbtProfileOutput].Project,
					Dataset:  profile.Outputs[formResponse.DbtProfileOutput].Dataset,
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
				connectionDetails = shared.ConnectionDetails{
					ConnType: formResponse.Warehouse,
					Username: formResponse.Username,
					Account:  formResponse.Account,
					Schema:   formResponse.Schema,
					Database: formResponse.Database,
				}
			}
		case "bigquery":
			{
				connectionDetails = shared.ConnectionDetails{
					ConnType: formResponse.Warehouse,
					Project:  formResponse.Project,
					Dataset:  formResponse.Dataset,
				}
			}
		default:
			{
				log.Fatalf("Unsupported connection type %v\n", formResponse.Warehouse)
			}
		}
	}
	return connectionDetails
}
